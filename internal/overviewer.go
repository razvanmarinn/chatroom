package main

import (
	"fmt"
	"log"
	"net/http"

	_db "github.com/razvanmarinn/chatroom/internal/db"

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

type RepositoryManager struct {
	RoomRepo    _db.RoomRepository
	MessageRepo _db.MessageRepository
	UserRepo    _db.UserRepository
}
type Overviewer struct {
	ConnectedClients  map[string][]*websocket.Conn
	RepositoryManager *RepositoryManager
}

func newOverviewer() *Overviewer {
	return &Overviewer{ConnectedClients: make(map[string][]*websocket.Conn), RepositoryManager: &RepositoryManager{}}
}

func (ow *Overviewer) connectWS(c echo.Context) error {
	ow.RepositoryManager.RoomRepo = c.Request().Context().Value("roomRepo").(_db.RoomRepository)
	ow.RepositoryManager.MessageRepo = c.Request().Context().Value("messageRepo").(_db.MessageRepository)
	ow.RepositoryManager.UserRepo = c.Request().Context().Value("userRepo").(_db.UserRepository)

	cookie, err := c.Cookie("session_token")
	if err != nil {

		fmt.Println("Error reading message:", err)
	}
	userUUID, found := ss.SessionStore.Get(cookie.Value)
	if !found {
		return c.String(http.StatusUnauthorized, "Invalid session")
	}
	roomName := c.Param("room_name")
	if !ow.RepositoryManager.RoomRepo.RoomExists(roomName) {
		return c.NoContent(400)
	}

	conn, err := websocketUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	ow.ConnectedClients[roomName] = append(ow.ConnectedClients[roomName], conn)

	go ow.handleMessages(conn, roomName, userUUID)

	return c.NoContent(200)
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

	room, err := ow.RepositoryManager.RoomRepo.GetRoomByName(roomName)
	if err != nil {
		fmt.Println("Error saving message:")

	}

	message, err := ow.RepositoryManager.MessageRepo.CreateMessage(room.ID, uuid.MustParse(userUUID), messageContent)
	if err != nil {
		fmt.Println("Error saving message:", err)
		return
	}

	user, err := ow.RepositoryManager.UserRepo.GetUserByID(uuid.MustParse(userUUID))
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
