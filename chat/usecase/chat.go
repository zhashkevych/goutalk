package usecase

import (
	"context"
	"crypto/sha1"
	"github.com/zhashkevych/goutalk/chat"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const salt = "gc7QRqMhWYHG7UgqpUbu"

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

	h := sha1.New()
	h.Write([]byte(password + salt))
	hashedPasswordBytes := h.Sum(nil)

	user = chat.NewUser(username, string(hashedPasswordBytes))

	if err := c.userRepo.Insert(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c *ChatEngine) GetAllUsers(ctx context.Context) ([]*chat.User, error) {
	return c.userRepo.GetAll(ctx)
}

func (c *ChatEngine) GetUserByID(ctx context.Context, id string) (*chat.User, error) {
	mid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return c.userRepo.GetByID(ctx, mid)
}

func (c *ChatEngine) CreateRoom(ctx context.Context, name string, creatorID string) error {
	return nil
}

func (c *ChatEngine) GetAllRooms(ctx context.Context) ([]*chat.Room, error) {
	return nil, nil
}

func (c *ChatEngine) GetRoomByID(ctx context.Context, id string) (*chat.Room, error) {
	return &chat.Room{}, nil
}

func (c *ChatEngine) AddRoomMember(ctx context.Context, roomID, memberID string) error {
	return nil
}

func (c *ChatEngine) RemoveRoomMeber(ctx context.Context, roomID, memberID string) error {
	return nil
}

func (c *ChatEngine) DeleteRoom(ctx context.Context, roomID string, user *chat.User) error {
	return nil
}
