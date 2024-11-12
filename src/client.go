package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Name string
	Id   uuid.UUID
	conn *websocket.Conn
}


func newClient(conn *websocket.Conn, name string) *Client {
	return &Client{Name: name, Id: uuid.New(), conn:conn}
}


func (c *Client) readMessage() {
	
}