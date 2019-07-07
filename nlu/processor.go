package nlu

import "context"

type Processor interface {
	Process(ctx context.Context, message, chatID string) (string, error)
}
