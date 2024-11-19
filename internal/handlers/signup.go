package handlers

import (
	"net/http"
	_db "github.com/razvanmarinn/chatroom/internal/db"

	"github.com/labstack/echo/v4"
)

func RegisterHandler(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "register_body", nil)
	}
	userRepo := c.Request().Context().Value("userRepo").(_db.UserRepository)
	username := c.FormValue("username")
	password := c.FormValue("password")

	if userRepo.UserExists(username) {

		return c.String(http.StatusConflict, "Username already taken")
	}

	_, err := userRepo.CreateUser(username, password)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to register user")
	}

	return c.Render(http.StatusOK, "index", nil)
}
