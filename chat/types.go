package chat

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}

func NewUser(username, pass string) *User {
	return &User{
		Username: username,
		Password: pass,
	}
}

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatorID string             `bson:"creator_id"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"created_at"`
	Members   []string           `bson:"members"`
}

func NewRoom(creatorID string, name string) *Room {
	return &Room{
		CreatorID: creatorID,
		CreatedAt: time.Now(),
		Name:      name,
		Members:   make([]string, 0),
	}
}
