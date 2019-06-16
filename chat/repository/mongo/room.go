package mongo

import (
	"context"
	"github.com/zhashkevych/goutalk/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": room.ID})
	return err
}

func (r *RoomsRepository) AddUser(ctx context.Context, room *chat.Room, u *chat.User) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": room.ID}, bson.M{"$addToSet": bson.M{"members": u.ID}})
	return err
}

func (r *RoomsRepository) RemoveUser(ctx context.Context, room *chat.Room, u *chat.User) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": room.ID}, bson.M{"$pull": bson.M{"members": u.ID}})
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
	res := r.db.FindOne(ctx, bson.M{"id": id})

	if err := res.Decode(&room); err != nil {
		return nil, err
	}

	return &room, nil
}

//func toRoom(record bson.M) *chat.Room {
//	var (
//		id        primitive.ObjectID
//		creatorID primitive.ObjectID
//		createdAt time.Time
//		members   []primitive.ObjectID
//	)
//
//	if _, ex := record["_id"]; ex {
//		id = record["_id"].(primitive.ObjectID)
//	}
//	if _, ex := record["creatorID"]; ex {
//		creatorID = record["creatorID"].(primitive.ObjectID)
//	}
//	if _, ex := record["createdAt"]; ex {
//		createdAt = record["createdAt"].(time.Time)
//	}
//	if _, ex := record["members"]; ex {
//		members = record["members"].([]primitive.ObjectID)
//	}
//
//	membersStr := make([]string, len(members))
//	for i := range members {
//		membersStr[i] = members[i].String()
//	}
//
//	return &chat.Room{
//		ID:        id.String(),
//		CreatedAt: createdAt,
//		CreatorID: creatorID.String(),
//		Members:   membersStr,
//	}
//}
