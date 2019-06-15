package chat

import (
	"github.com/satori/go.uuid"
	"time"
)

type User struct {
	ID       uuid.UUID `json:"id" bson:"id"`
	Username string    `json:"username" bson:"username"`
	Password string    `json:"password" bson:"password"`
}

func NewUser(username, pass string) *User {
	return &User{
		ID:       uuid.NewV4(),
		Username: username,
		Password: pass,
	}
}

type Room struct {
	ID        uuid.UUID    `json:"id" bson:"id"`
	CreatorID uuid.UUID    `json:"creatorID" bson:"creatorID"`
	CreatedAt time.Time    `json:"createdAt" bson:"createdAt"`
	Members   []uuid.UUID `json:"members" bson:"members"`
}

func NewRoom(creatorID uuid.UUID) *Room {
	return &Room{
		ID:        uuid.NewV4(),
		CreatorID: creatorID,
		CreatedAt: time.Now(),
		Members:   make([]uuid.UUID, 0),
	}
}
