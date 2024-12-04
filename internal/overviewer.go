package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/razvanmarinn/chatroom/internal/services"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	ss "github.com/razvanmarinn/chatroom/internal/session_store"
)

var websocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const _NO_OF_CACHED_MESSAGES = 100

type Overviewer struct {
	ConnectedClients map[string][]*websocket.Conn
	ServiceManager   *services.ServiceManager
}

func newOverviewer() *Overviewer {
	return &Overviewer{ConnectedClients: make(map[string][]*websocket.Conn), ServiceManager: &services.ServiceManager{}}
}

func (ow *Overviewer) connectWS(c echo.Context) error {
	ow.ServiceManager = c.Request().Context().Value("serviceManager").(*services.ServiceManager)

	cookie, err := c.Cookie("session_token")
	if err != nil {
		fmt.Println("Error reading session token:", err)
		return c.String(http.StatusUnauthorized, "Invalid session")
	}

	userUUID, found := ss.SessionStore.Get(cookie.Value)
	if !found {
		return c.String(http.StatusUnauthorized, "Invalid session")
	}

	roomName := c.Param("room_name")
	if !ow.ServiceManager.RoomService.RoomExists(roomName) {
		return c.NoContent(http.StatusBadRequest)
	}

	conn, err := websocketUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	room, err := ow.ServiceManager.RoomService.GetRoomByName(roomName)
	if err != nil {
		fmt.Println("Error getting room:", err)
		conn.Close()
		return c.String(http.StatusInternalServerError, "Room not found")
	}

	messages, err := ow.ServiceManager.MessageService.GetLastMessagesByRoomID(room.ID, _NO_OF_CACHED_MESSAGES)
	if err != nil {
		fmt.Println("Error fetching room messages:", err)

	}

	for _, message := range messages {
		user, err := ow.ServiceManager.UserService.GetUserByID(message.UserID)
		if err != nil {
			fmt.Println("Error getting user for message:", err)
			continue
		}

		formattedMessage := fmt.Sprintf("%s: %s", user.Username, message.Content)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(formattedMessage)); err != nil {
			fmt.Println("Error sending historical message:", err)
			break
		}
	}

	ow.ConnectedClients[roomName] = append(ow.ConnectedClients[roomName], conn)

	go ow.handleMessages(conn, roomName, userUUID)

	return c.NoContent(http.StatusOK)
}

func (ow *Overviewer) handleMessages(conn *websocket.Conn, roomName string, userUUID string) {
	defer conn.Close()

	defer ow.removeConnectionFromChatroom(conn, roomName)

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {

			fmt.Println("Error reading message:", err)
			break
		}
		ow.broadcastMessageToChatroom(roomName, msg, userUUID)
	}
}

func (ow *Overviewer) removeConnectionFromChatroom(conn *websocket.Conn, roomName string) {
	connections := ow.ConnectedClients[roomName]
	for i, c := range connections {
		if c == conn {
			ow.ConnectedClients[roomName] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}

func (ow *Overviewer) broadcastMessageToChatroom(roomName string, messageContent []byte, userUUID string) {

	room, err := ow.ServiceManager.RoomService.GetRoomByName(roomName)
	if err != nil {
		fmt.Println("Error saving message:")

	}

	message, err := ow.ServiceManager.MessageService.CreateMessage(room.ID, uuid.MustParse(userUUID), messageContent)
	if err != nil {
		fmt.Println("Error saving message:", err)
		return
	}

	user, err := ow.ServiceManager.UserService.GetUserByID(uuid.MustParse(userUUID))
	if err != nil {
		fmt.Println("Error saving message:")

	}

	formattedMessage := fmt.Sprintf("%s: %s", user.Username, message.Content)

	connections := ow.ConnectedClients[roomName]
	for _, conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(formattedMessage)); err != nil {
			fmt.Println("Error sending message:", err)
			continue
		}
	}
}
