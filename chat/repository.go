package chat

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Insert(ctx context.Context, u *User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetAll(ctx context.Context) ([]*User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type RoomRepository interface {
	Insert(ctx context.Context, r *Room) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetAll(ctx context.Context) ([]*Room, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*Room, error)
	AddUser(ctx context.Context, roomID, userID primitive.ObjectID) error
	RemoveUser(ctx context.Context, roomID, userID primitive.ObjectID) error
}
