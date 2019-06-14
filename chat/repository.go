package chat

import "github.com/satori/go.uuid"

type UserRepository interface {
	Create(username, password string) (string, error)
	Delete(id uuid.UUID)
}

type RoomRepository interface {
	Create()
	Delete()
	AddUser()
	RemoveUser()
}