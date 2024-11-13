package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/razvanmarinn/chatroom/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	var user db.User
	result := db.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {

		return c.String(http.StatusUnauthorized, "Invalid username or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {

		return c.String(http.StatusUnauthorized, "Invalid username or password")
	}

	return c.Render(http.StatusOK, "login_body", nil)
}
