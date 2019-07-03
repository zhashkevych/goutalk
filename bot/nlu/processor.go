package nlu

import "context"

type Processor interface {
	Process(ctx context.Context, message, userID, chatID string) (string, error)
}
