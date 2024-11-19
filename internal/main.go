package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/razvanmarinn/chatroom/internal/db"
	r_fact "github.com/razvanmarinn/chatroom/internal/db/repository_factory"
	"github.com/razvanmarinn/chatroom/internal/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/razvanmarinn/chatroom/internal/cfg"
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

	config := cfg.LoadConfig()

	dbConn, err := db.InitDatabase(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	repoFactory, err := r_fact.CreateRepositoryFactory(config.DbType, dbConn)
	if err != nil {
		log.Fatalf("Failed to create repository factory: %v", err)
	}

	userRepo, err := repoFactory.CreateUserRepository()
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	roomRepo, err := repoFactory.CreateRoomRepository()
	if err != nil {
		log.Fatalf("Failed to create room repository: %v", err)
	}

	messageRepo, err := repoFactory.CreateMessageRepository()
	if err != nil {
		log.Fatalf("Failed to create message repository: %v", err)
	}
	//TODO: Add services 

	e := echo.New()
	godotenv.Load("../.env")

	ow := newOverviewer()

	t := &Template{
		templates: template.Must(template.ParseGlob("frontend/*.html")),
	}
	e.Renderer = t
	e.Use(middleware.AddRepositoriesToContext(&userRepo, &roomRepo, &messageRepo))
	e.GET("/login", handlers.LoginHandler)
	e.POST("/login", handlers.LoginHandler)
	e.GET("/signup", handlers.RegisterHandler)
	e.POST("/signup", handlers.RegisterHandler)

	protected := e.Group("")
	protected.Use(IsAuthenticated)

	protected.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})
	protected.POST("/create-room", handlers.RoomHandler)
	protected.GET("/ws/room/:room_name", ow.connectWS)
	log.Fatal(e.Start(":8080"))
}
