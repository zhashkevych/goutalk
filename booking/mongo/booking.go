package mongo

import (
	"context"
	"github.com/zhashkevych/goutalk/booking"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	bookingsCollection = "bookings"
)

//var ErrorNotFound = errors.New("item not found")

type BookingRepository struct {
	db *mongo.Collection
}

func NewBookingRepository(db *mongo.Database) *BookingRepository {
	return &BookingRepository{
		db: db.Collection(bookingsCollection),
	}
}

func (r *BookingRepository) Insert(ctx context.Context, item *booking.BookItem) error {
	res, err := r.db.InsertOne(ctx, item)
	if err != nil {
		return err
	}

	item.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *BookingRepository) Delete(ctx context.Context) error {
	//_, err := r.db.DeleteOne(ctx, bson.M{"_id": id})
	//if err != nil {
	//	return chat.NewErrorNotFound("room", "id", id.Hex())
	//}

	return nil
}

func (r *BookingRepository) Update(ctx context.Context, id, date string) error {
	pID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.db.UpdateOne(ctx, bson.M{"_id": pID}, bson.M{"$set": bson.M{"date": date}})
	return err
}

func (r *BookingRepository) GetByUserID(ctx context.Context, userID string) ([]*booking.BookItem, error) {
	cur, err := r.db.Find(ctx, bson.M{"user_id": userID})
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	out := make([]*booking.BookItem, 0)

	for cur.Next(ctx) {
		var item booking.BookItem
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}

		out = append(out, &item)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *BookingRepository) RemoveByUserID(ctx context.Context, userID string) error {
	_, err := r.db.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}

	return nil
}

func (r *BookingRepository) GetAll(ctx context.Context) ([]*booking.BookItem, error) {
	//cur, err := r.db.Find(ctx, bson.D{})
	//defer cur.Close(ctx)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	out := make([]*booking.BookItem, 0)
	//
	//for cur.Next(ctx) {
	//	var room chat.Room
	//	err := cur.Decode(&room)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	out = append(out, &room)
	//}
	//if err := cur.Err(); err != nil {
	//	return nil, err
	//}

	return out, nil
}
