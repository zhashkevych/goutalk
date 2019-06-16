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
	// getting hash from password
	h := sha1.New()
	h.Write([]byte(password + salt))
	hashedPassword := string(h.Sum(nil))

	// getting user from db and comparing password hashes
	user, err := c.userRepo.GetByUsername(ctx, username)
	if err == nil {
		if user.Password == hashedPassword {
			return user, nil
		}

		return nil, chat.ErrWrongPassword
	}

	// creating new user
	user = chat.NewUser(username, hashedPassword)

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

func (c *ChatEngine) CreateRoom(ctx context.Context, name string, creatorID string) (*chat.Room, error) {
	cid, err := primitive.ObjectIDFromHex(creatorID)
	if err != nil {
		return nil, err
	}

	room := chat.NewRoom(cid, name)

	return room, c.roomRepo.Insert(ctx, room)
}

func (c *ChatEngine) GetAllRooms(ctx context.Context) ([]*chat.Room, error) {
	return c.roomRepo.GetAll(ctx)
}

func (c *ChatEngine) GetRoomByID(ctx context.Context, id string) (*chat.Room, error) {
	mid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return c.roomRepo.GetByID(ctx, mid)
}

func (c *ChatEngine) AddRoomMember(ctx context.Context, roomID, memberID string) error {
	rID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}

	mID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return err
	}

	if _, err := c.userRepo.GetByID(ctx, mID); err != nil {
		return chat.NewErrorNotFound("user", "id", memberID)
	}

	return c.roomRepo.AddUser(ctx, rID, mID)
}

func (c *ChatEngine) RemoveRoomMeber(ctx context.Context, roomID, memberID string) error {
	rID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}

	mID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return err
	}

	return c.roomRepo.RemoveUser(ctx, rID, mID)
}

func (c *ChatEngine) DeleteRoom(ctx context.Context, roomID string, user *chat.User) error {
	rID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}

	room, err := c.roomRepo.GetByID(ctx, rID)
	if err != nil {
		return err
	}

	if room.CreatorID != user.ID {
		return chat.ErrMissingAccessRights
	}

	return c.roomRepo.Delete(ctx, rID)
}

func (c *ChatEngine) GetRoomMembers(ctx context.Context, roomID string) ([]*chat.User, error) {
	rID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, err
	}

	room, err := c.roomRepo.GetByID(ctx, rID)
	if err != nil {
		return nil, err
	}

	users, err := c.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*chat.User, len(room.Members))
	for i := range room.Members {
		for _, user := range users {
			if user.ID == room.Members[i] {
				out[i] = user
			}
		}
	}

	return out, nil
}
