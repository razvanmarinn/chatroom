package services

type ServiceManager struct {
	UserService    *UserService
	RoomService    *RoomService
	MessageService *MessageService
}
