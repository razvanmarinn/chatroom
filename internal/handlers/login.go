package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/razvanmarinn/chatroom/internal/db"
	ss "github.com/razvanmarinn/chatroom/internal/session_store"

	"golang.org/x/crypto/bcrypt"
)

func generateSecureSessionToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err) //
	}
	return hex.EncodeToString(bytes)
}

func LoginHandler(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "login_body", nil)
	}

	userRepo := c.Request().Context().Value("userRepo").(db.UserRepository)

	username := c.FormValue("username")
	password := c.FormValue("password")

	if !userRepo.UserExists(username) {
		return c.String(http.StatusUnauthorized, "Invalid username or password")
	}
	user, err := userRepo.GetUserByUsername(username)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid username or password")

	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {

		return c.String(http.StatusUnauthorized, "Invalid username or password")
	}
	sessionToken := generateSecureSessionToken()
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true, // Prevent JavaScript access
		Secure:   true, // Use HTTPS
	})

	ss.SessionStore.Set(sessionToken, user.ID.String())

	return c.Redirect(http.StatusSeeOther, "/")
}
