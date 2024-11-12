package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var globalCounter = 0

type User struct {
	Name string
	Id   int
}

func newUser(name string) *User {
	globalCounter++
	return &User{Name: name, Id: globalCounter}
}

type Page struct {
	Data []User
}

func newPage() *Page {
	return &Page{
		Data: []User{},
	}
}

func main() {
	e := echo.New()

	ow := newOverviewer()

	page := newPage()
	t := &Template{
		templates: template.Must(template.ParseFiles("index.html")),
	}
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", page)
	})

	e.PUT("/change_chatroom", func(c echo.Context) error {
		fmt.Println("chat changed to room no : ", c.FormValue("chatroom_id"))
		return c.Render(http.StatusOK, "chat_input", c.FormValue("chatroom_id"))
	})

	e.GET("/ws/:chatroom_id", ow.connectWS)

	e.POST("/create_user", func(c echo.Context) error {
		name := c.FormValue("name")
		user := newUser(name)
		page.Data = append(page.Data, *user)
		return c.Render(http.StatusOK, "displaying", page)
	})

	log.Fatal(e.Start(":8080"))
}
