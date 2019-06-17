package mongo

import (
	"context"
	"github.com/zhashkevych/goutalk/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	roomsCollection = "rooms"
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
	res, err := r.db.InsertOne(ctx, room)
	if err != nil {
		return err
	}

	room.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *RoomsRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return chat.NewErrorNotFound("room", "id", id.Hex())
	}

	return nil
}

func (r *RoomsRepository) AddUser(ctx context.Context, roomID, userID primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": roomID}, bson.M{"$addToSet": bson.M{"members": userID.Hex()}})
	return err
}

func (r *RoomsRepository) RemoveUser(ctx context.Context, roomID, userID primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": roomID}, bson.M{"$pull": bson.M{"members": userID.Hex()}})
	return err
}

func (r *RoomsRepository) GetAll(ctx context.Context) ([]*chat.Room, error) {
	cur, err := r.db.Find(ctx, bson.D{})
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	out := make([]*chat.Room, 0)

	for cur.Next(ctx) {
		var room chat.Room
		err := cur.Decode(&room)
		if err != nil {
			return nil, err
		}

		out = append(out, &room)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *RoomsRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*chat.Room, error) {
	var room chat.Room
	res := r.db.FindOne(ctx, bson.M{"_id": id})

	if err := res.Decode(&room); err != nil {
		return nil, chat.NewErrorNotFound("room", "id", id.Hex())
	}

	return &room, nil
}