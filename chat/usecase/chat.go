package usecase

import (
	"context"
	"github.com/satori/go.uuid"
	"github.com/zhashkevych/goutalk/chat"
)

type ChatEngine struct {
	userRepo chat.UserRepository
	roomRepo chat.RoomRepository
}

func NewChatEngine(userRepo chat.UserRepository, roomRepo chat.RoomRepository) *ChatEngine {
	return &ChatEngine{
		userRepo: userRepo,
		roomRepo: roomRepo,
	}
}

func (c *ChatEngine) RegisterNewUser(ctx context.Context, name, nickname, password string) error {
	// TODO: check if there is such user

	u := chat.NewUser(name, nickname, password)

}

func (c *ChatEngine) GetAllUsers(ctx context.Context) ([]*chat.User, error) {
	return nil, nil
}

func (c *ChatEngine) GetUserByID(ctx context.Context, id uuid.UUID) (*chat.User, error) {
	return &chat.User{}, nil
}

func (c *ChatEngine) CreateRoom(ctx context.Context, name string, creatorID uuid.UUID) error {
	return nil
}

func (c *ChatEngine) GetAllRooms(ctx context.Context) ([]*chat.Room, error) {
	return nil, nil
}

func (c *ChatEngine) GetRoomByID(ctx context.Context, id uuid.UUID) (*chat.Room, error) {
	return &chat.Room{}, nil
}

func (c *ChatEngine) AddRoomMember(ctx context.Context, roomID, memberID uuid.UUID) error {
	return nil
}

func (c *ChatEngine) RemoveRoomMeber(ctx context.Context, roomID, memberID uuid.UUID) error {
	return nil
}

func (c *ChatEngine) DeleteRoom(ctx context.Context, roomID, user *chat.User) error {
	return nil
}
