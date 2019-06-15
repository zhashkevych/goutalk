package mongo

import (
	"context"
	"github.com/satori/go.uuid"
	"github.com/zhashkevych/goutalk/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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

func (r *RoomsRepository) GetAll(ctx context.Context) ([]*chat.Room, error) {
	cur, err := r.db.Find(ctx, bson.D{})
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	out := make([]*chat.Room, 0)

	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}

		out = append(out, toRoom(result))
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *RoomsRepository) GetByID(ctx context.Context, id uuid.UUID) (*chat.Room, error) {
	var record bson.M
	res := r.db.FindOne(ctx, bson.M{"id": id})

	if err := res.Decode(record); err != nil {
		return nil, err
	}

	return toRoom(record), nil
}

func toRoom(record bson.M) *chat.Room {
	var (
		id        uuid.UUID
		creatorID uuid.UUID
		createdAt time.Time
		members   []uuid.UUID
	)

	if _, ex := record["id"]; ex {
		id = record["id"].(uuid.UUID)
	}
	if _, ex := record["creatorID"]; ex {
		creatorID = record["creatorID"].(uuid.UUID)
	}
	if _, ex := record["createdAt"]; ex {
		createdAt = record["createdAt"].(time.Time)
	}
	if _, ex := record["members"]; ex {
		members = record["members"].([]uuid.UUID)
	}

	return &chat.Room{
		ID:        id,
		CreatedAt: createdAt,
		CreatorID: creatorID,
		Members:   members,
	}
}
