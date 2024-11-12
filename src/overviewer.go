package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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
	chatroomID := c.Param("chatroom_id")
	conn, err := websocketUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	// client := newClient(conn, "test")
	chatroomConnections[chatroomID] = append(chatroomConnections[chatroomID], conn)

	go handleMessages(conn, chatroomID)

	return c.NoContent(200)
}

func handleMessages(conn *websocket.Conn, chatroomID string) {
	defer conn.Close()

	defer removeConnectionFromChatroom(conn, chatroomID)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		broadcastMessageToChatroom(chatroomID, msg)
	}
}

func removeConnectionFromChatroom(conn *websocket.Conn, chatroomID string) {
	connections := chatroomConnections[chatroomID]
	for i, c := range connections {
		if c == conn {
			chatroomConnections[chatroomID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}

func broadcastMessageToChatroom(chatroomID string, message []byte) {
	connections := chatroomConnections[chatroomID]

	for _, conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error sending message:", err)
			continue
		}
	}
}
