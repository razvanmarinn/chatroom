package main

import "github.com/google/uuid"

var MAX_CLIENTS_PER_CHATROOM = 256

type ChatRoom struct {
	ActiveConnectionsCount int
	Name              string
	Id              uuid.UUID
}