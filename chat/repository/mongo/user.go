package mongo

import (
	"context"
	"github.com/zhashkevych/goutalk/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	res, err := r.db.InsertOne(ctx, u)
	if err != nil {
		return err
	}

	u.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return chat.NewErrorNotFound("user", "id", id.Hex())
	}

	return nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*chat.User, error) {
	cur, err := r.db.Find(ctx, bson.D{})
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	out := make([]*chat.User, 0)

	for cur.Next(ctx) {
		var user chat.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}

		out = append(out, &user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*chat.User, error) {
	var user chat.User
	res := r.db.FindOne(ctx, bson.M{"_id": id})

	if err := res.Decode(&user); err != nil {
		return nil, chat.NewErrorNotFound("user", "id", id.Hex())
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*chat.User, error) {
	var user chat.User
	res := r.db.FindOne(ctx, bson.M{"username": username})

	if err := res.Decode(&user); err != nil {
		return nil, chat.NewErrorNotFound("user", "username", username)
	}

	return &user, nil
}
