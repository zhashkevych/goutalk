package booking

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookItem struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        string             `bson:"user_id,omitempty"`
	Date          string             `bson:"date,omitempty"`
	Location      string             `bson:"location,omitempty"`
	VenueType     string             `bson:"venue_type,omitempty"`
	VenueTitle    string             `bson:"venue_title,omitempty"`
	VenueFacility string             `bson:"venue_facility,omitempty"`
}

type Repository interface {
	Insert(ctx context.Context, item *BookItem) error
	Delete(ctx context.Context) error
	Update(ctx context.Context, id, date string) error
	GetByUserID(ctx context.Context, userID string) ([]*BookItem, error)
	GetByID(ctx context.Context, id string) (*BookItem, error)
	RemoveByUserID(ctx context.Context, userID string) error
	RemoveByID(ctx context.Context, userID, id string) error
	GetAll(ctx context.Context) ([]*BookItem, error)
}
