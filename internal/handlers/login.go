package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/razvanmarinn/chatroom/internal/logger"
	"github.com/razvanmarinn/chatroom/internal/services"
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

	userService := c.Request().Context().Value("serviceManager").(*services.ServiceManager).UserService
	logger := c.Request().Context().Value("logger").(logger.Logger)


	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.String(http.StatusBadRequest, "Username and password are required")
	}

	user, err := userService.GetUserByUsername(username)
	if err != nil {
		logger.Error("Login attempt failed for user %s: %v", username, err)
		return c.String(http.StatusUnauthorized, "Invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
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
