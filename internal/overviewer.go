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
	"gorm.io/gorm"
)

var chatroomConnections = make(map[string][]*websocket.Conn)

var websocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Overviewer struct {
	ConnectedClients []Client
}

func newOverviewer() *Overviewer {
	return &Overviewer{ConnectedClients: []Client{}}
}

func (ow *Overviewer) connectWS(c echo.Context) error {
	cookie, err := c.Cookie("session_token")
	userUUID, found := ss.SessionStore.Get(cookie.Value)
	if !found {
		return c.String(http.StatusUnauthorized, "Invalid session")
	}
	roomName := c.Param("room_name")
	if !_db.CheckRoomExists(roomName) {
		return c.NoContent(400)
	}

	conn, err := websocketUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	chatroomConnections[roomName] = append(chatroomConnections[roomName], conn)

	go handleMessages(conn, roomName, userUUID)

	return c.NoContent(200)
}

func handleMessages(conn *websocket.Conn, roomName string, userUUID string) {
	defer conn.Close()

	defer removeConnectionFromChatroom(conn, roomName)

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {

			fmt.Println("Error reading message:", err)
			break
		}
		broadcastMessageToChatroom(_db.DB, roomName, msg, userUUID)
	}
}

func removeConnectionFromChatroom(conn *websocket.Conn, roomName string) {
	connections := chatroomConnections[roomName]
	for i, c := range connections {
		if c == conn {
			chatroomConnections[roomName] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}

func broadcastMessageToChatroom(db *gorm.DB, roomName string, messageContent []byte, userUUID string) {

	roomUUID, err := _db.GetRoomUUID(roomName)
	if err != nil {
		fmt.Println("Error saving message:")

	}
	message := &_db.Message{
		RoomID:  roomUUID,
		UserID:  uuid.MustParse(userUUID),
		Content: string(messageContent),
	}

	result := db.Create(message)
	if result.Error != nil {
		fmt.Println("Error saving message:", result.Error)
		return
	}

	username, err :=_db.GetUsername(uuid.MustParse(userUUID))
	if err != nil {
		fmt.Println("Error saving message:")

	}

	formattedMessage := fmt.Sprintf("%s: %s", username, string(messageContent))

	connections := chatroomConnections[roomName]
	for _, conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(formattedMessage)); err != nil {
			fmt.Println("Error sending message:", err)
			continue
		}
	}
}
