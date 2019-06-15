package mongo

import (
	"context"
	"github.com/zhashkevych/goutalk/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	usersCollection = "users"
)

type UserRepository struct {
	db 		*mongo.Collection
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
