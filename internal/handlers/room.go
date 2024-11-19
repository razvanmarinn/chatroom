package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	ss "github.com/razvanmarinn/chatroom/internal/session_store"

	"github.com/razvanmarinn/chatroom/internal/db"
)

func RoomHandler(c echo.Context) error {
	cookie, err := c.Cookie("session_token")
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	userUUID, found := ss.SessionStore.Get(cookie.Value)
	if !found {
		return c.String(http.StatusUnauthorized, "Invalid session")
	}

	roomName := c.FormValue("cr_room_name")

	var existingRoom db.Room
	if err := db.DB.Where("room_name = ?", roomName).First(&existingRoom).Error; err == nil {
		return c.String(http.StatusConflict, "Room name already exists")
	}
	owner, err := uuid.Parse(userUUID)
	if err != nil {
		return err
	}

	room := &db.Room{
		ID:       uuid.New(),
		RoomName: roomName,
		Owner:    owner,
	}

	if err := db.DB.Create(&room).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create room")
	}

	return c.HTML(http.StatusOK, "<p>Room successfully created!</p>")
}
