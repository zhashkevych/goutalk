package chat

import (
	"context"
)

type UserRepository interface {
	// TODO: extend with GET methods
	Insert(ctx context.Context, u *User) error
	Delete(ctx context.Context, u *User) error
}

type RoomRepository interface {
	Insert(ctx context.Context, r *Room) error
	Delete(ctx context.Context, r *Room) error
	AddUser(ctx context.Context, r *Room, u *User) error
	RemoveUser(ctx context.Context, r *Room, u *User) error
}
