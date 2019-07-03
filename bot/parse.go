package bot

import (
	"encoding/json"
	"strings"
)

const mention = "@bot"

type Message struct {
	Text   string `json:"message"`
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

func Parse(msg []byte) (*Message, error) {
	m := new(Message)
	if err := json.Unmarshal(msg, m); err != nil {
		return nil, err
	}

	if !strings.Contains(m.Text, mention) {
		return nil, nil
	}

	m.Text = strings.TrimPrefix(m.Text, mention + " ")

	return m, nil
}
