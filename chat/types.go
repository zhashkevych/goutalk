package chat

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

func NewUser(username, pass string) *User {
	return &User{
		Username: username,
		Password: pass,
	}
}

type Room struct {
	ID        primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	CreatorID primitive.ObjectID    `json:"creatorID" bson:"creatorID"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Members   []primitive.ObjectID  `json:"members" bson:"members"`
}

func NewRoom(creatorID primitive.ObjectID) *Room {
	return &Room{
		CreatorID: creatorID,
		CreatedAt: time.Now(),
		Members:   make([]primitive.ObjectID, 0),
	}
}
