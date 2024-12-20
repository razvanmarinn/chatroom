package handlers

import (
	"net/http"

	"github.com/razvanmarinn/chatroom/internal/logger"
	"github.com/razvanmarinn/chatroom/internal/services"

	"github.com/labstack/echo/v4"
)

func RegisterHandler(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "register_body", nil)
	}
	userService := c.Request().Context().Value("serviceManager").(*services.ServiceManager).UserService
	logger := c.Request().Context().Value("logger").(logger.Logger)

	username := c.FormValue("username")
	password := c.FormValue("password")

	if userService.CheckUserExists(username) {
		logger.Error("Username already taken " + username)
		return c.String(http.StatusConflict, "Username already taken")
	}

	_, err := userService.CreateUser(username, password)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to register user")
	}

	return c.Render(http.StatusOK, "index", nil)
}
