package mongo

import (
	"context"
	"github.com/satori/go.uuid"
	"github.com/zhashkevych/goutalk/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	usersCollection = "users"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db: db.Collection(usersCollection),
	}
}

func (r *UserRepository) Insert(ctx context.Context, u *chat.User) error {
	_, err := r.db.InsertOne(ctx, u)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, u *chat.User) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"id": u.ID})
	return err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*chat.User, error) {
	cur, err := r.db.Find(ctx, bson.D{})
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	out := make([]*chat.User, 0)

	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}

		out = append(out, toUser(result))
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*chat.User, error) {
	var record bson.M
	res := r.db.FindOne(ctx, bson.M{"id": id})

	if err := res.Decode(record); err != nil {
		return nil, err
	}

	return toUser(record), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*chat.User, error) {
	var record bson.M
	res := r.db.FindOne(ctx, bson.M{"username": username})

	if err := res.Decode(record); err != nil {
		return nil, err
	}

	return toUser(record), nil
}

func toUser(record bson.M) *chat.User {
	var (
		username string
		id       uuid.UUID
	)

	if _, ex := record["id"]; ex {
		id = record["id"].(uuid.UUID)
	}
	if _, ex := record["username"]; ex {
		username = record["username"].(string)
	}

	return &chat.User{
		ID:       id,
		Username: username,
	}
}
