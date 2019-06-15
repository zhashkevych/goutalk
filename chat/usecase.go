package chat

import (
	"context"
	"github.com/satori/go.uuid"
)

type UseCase interface {
	LoginUser(ctx context.Context, username, password string) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)

	CreateRoom(ctx context.Context, name string, creatorID uuid.UUID) error
	GetAllRooms(ctx context.Context) ([]*Room, error)
	GetRoomByID(ctx context.Context, id uuid.UUID) (*Room, error)
	AddRoomMember(ctx context.Context, roomID, memberID uuid.UUID) error
	RemoveRoomMeber(ctx context.Context, roomID, memberID uuid.UUID) error
	DeleteRoom(ctx context.Context, roomID, user *User) error
}