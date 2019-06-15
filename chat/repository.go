package chat

import (
	"context"
	"github.com/satori/go.uuid"
)

type UserRepository interface {
	Insert(ctx context.Context, u *User) error
	Delete(ctx context.Context, u *User) error
	GetAll(ctx context.Context) ([]*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type RoomRepository interface {
	Insert(ctx context.Context, r *Room) error
	Delete(ctx context.Context, r *Room) error
	GetAll(ctx context.Context) ([]*Room, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Room, error)
	AddUser(ctx context.Context, r *Room, u *User) error
	RemoveUser(ctx context.Context, r *Room, u *User) error
}
