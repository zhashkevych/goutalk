package chat

import (
	"github.com/satori/go.uuid"
	"time"
)

type User struct {
	ID       uuid.UUID `json:"id" bson:"id"`
	Name     string    `json:"name" bson:"name"`
	Nickname string    `json:"nickname" bson:"nickname"`
	Password string    `json:"password" bson:"password"`
}

func NewUser(name, nickname, pass string) *User {
	return &User{
		ID:       uuid.NewV4(),
		Name:     name,
		Nickname: nickname,
		Password: pass,
	}
}

type Room struct {
	ID        uuid.UUID `json:"id" bson:"id"`
	CreatorID uuid.UUID `json:"creatorID" bson:"creatorID"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Members   []*User   `json:"members" bson:"members"`
}

func NewRoom(creatorID uuid.UUID) *Room {
	return &Room{
		ID:        uuid.NewV4(),
		CreatorID: creatorID,
		CreatedAt: time.Now(),
		Members:   make([]*User, 0),
	}
}
