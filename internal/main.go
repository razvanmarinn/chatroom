package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/razvanmarinn/chatroom/internal/auth"
	"github.com/razvanmarinn/chatroom/internal/db"
	"github.com/razvanmarinn/chatroom/internal/handlers"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			return c.Redirect(http.StatusFound, "/login")
		}
		if sessionToken == nil || sessionToken.Value == "" {
			return c.Redirect(http.StatusFound, "/login")
		}

		return next(c)
	}
}
func main() {
	e := echo.New()
	godotenv.Load("../.env")
	db.Init()
	ow := newOverviewer()

	t := &Template{
		templates: template.Must(template.ParseGlob("frontend/*.html")),
	}
	e.Renderer = t

	e.GET("/login", auth.LoginHandler)
	e.POST("/login", auth.LoginHandler)
	e.GET("/signup", auth.RegisterHandler)
	e.POST("/signup", auth.RegisterHandler)

	protected := e.Group("")
	protected.Use(IsAuthenticated)

	protected.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})
	protected.POST("/create-room", handlers.RoomHandler)
	protected.GET("/ws/room/:room_name", ow.connectWS)
	log.Fatal(e.Start(":8080"))
}
