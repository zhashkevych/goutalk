package chat

import (
	"context"
)

type UseCase interface {
	LoginUser(ctx context.Context, username, password string) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)

	CreateRoom(ctx context.Context, name string, creatorID string) (*Room, error)
	GetAllRooms(ctx context.Context) ([]*Room, error)
	GetRoomByID(ctx context.Context, id string) (*Room, error)
	AddRoomMember(ctx context.Context, roomID, memberID string) error
	RemoveRoomMeber(ctx context.Context, roomID, memberID string) error
	DeleteRoom(ctx context.Context, roomID string, user *User) error
	GetRoomMembers(ctx context.Context, roomID string) ([]*User, error)

	SendMessage(message *Message) error
}
