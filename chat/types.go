package chat

import (
	"github.com/satori/go.uuid"
	"time"
)

type User struct {
	ID       uuid.UUID
	Name     string
	Nickname string
	Password string
}

type Room struct {
	ID        uuid.UUID
	CreatorID uuid.UUID
	CreatedAt time.Time
	Members   []*User
}
