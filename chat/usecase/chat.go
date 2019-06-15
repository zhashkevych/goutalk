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

func (c *ChatEngine) LoginUser(ctx context.Context, username, password string) (*chat.User, error) {
	user, err := c.userRepo.GetByUsername(ctx, username)
	if err == nil {
		return user, nil
	}

	user = chat.NewUser(username, password)

	if err := c.userRepo.Insert(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
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
