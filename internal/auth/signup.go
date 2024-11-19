package auth

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/razvanmarinn/chatroom/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func RegisterHandler(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "register_body", nil)
	}
	username := c.FormValue("username")
	password := c.FormValue("password")

	var existingUser db.User
	result := db.DB.Where("username = ?", username).First(&existingUser)

	if result.Error == nil {

		return c.String(http.StatusConflict, "Username already taken")
	}

	encryptedPassword, err := HashPassword(password)
	if err != nil {

		return c.String(http.StatusInternalServerError, "Failed to hash password")
	}

	newUser := db.User{
		ID:       uuid.New(),
		Username: username,
		Password: encryptedPassword,
	}

	if err := db.DB.Create(&newUser).Error; err != nil {

		return c.String(http.StatusInternalServerError, "Failed to register user")
	}

	return c.Render(http.StatusOK, "index", nil)
}
