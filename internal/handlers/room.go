package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	ss "github.com/razvanmarinn/chatroom/internal/session_store"

	"github.com/razvanmarinn/chatroom/internal/db"
)

func RoomHandler(c echo.Context) error {
	roomRepo := c.Request().Context().Value("roomRepo").(db.RoomRepository)

	cookie, err := c.Cookie("session_token")
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	userUUID, found := ss.SessionStore.Get(cookie.Value)
	if !found {
		return c.String(http.StatusUnauthorized, "Invalid session")
	}

	roomName := c.FormValue("cr_room_name")

	if roomRepo.RoomExists(roomName) {
		return c.String(http.StatusConflict, "Room name already exists")
	}

	owner, err := uuid.Parse(userUUID)
	if err != nil {
		return err
	}

	room_id, err := roomRepo.CreateRoom(roomName, owner)
	if err != nil {
		return c.String(http.StatusConflict, "Room name already exists") // TODO: improve error
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>Room successfully created with id %s</p>", room_id.String()))
}
