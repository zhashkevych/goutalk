package bot

import "context"

type LanguageProcessor interface {
	Process(ctx context.Context, message, userID, chatID string) error
}
