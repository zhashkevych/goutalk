package mongo

import (
	"context"
	"github.com/zhashkevych/goutalk/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	roomsCollection = "users"
)

type RoomsRepository struct {
	db *mongo.Collection
}

func NewRoomsRepository(db *mongo.Database) *RoomsRepository {
	return &RoomsRepository{
		db: db.Collection(roomsCollection),
	}
}

func (r *RoomsRepository) Insert(ctx context.Context, room *chat.Room) error {
	_, err := r.db.InsertOne(ctx, r)
	return err
}

func (r *RoomsRepository) Delete(ctx context.Context, room *chat.Room) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"id": room.ID})
	return err
}

func (r *RoomsRepository) AddUser(ctx context.Context, room *chat.Room, u *chat.User) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"id": room.ID}, bson.M{"$addToSet": bson.M{"members": u.ID}})
	return err
}

func (r *RoomsRepository) RemoveUser(ctx context.Context, room *chat.Room, u *chat.User) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"id": room.ID}, bson.M{"$pull": bson.M{"members": u.ID}})
	return err
}
