package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"new_todo_project/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserDb interface {
	Add(ctx context.Context, user domain.User) (string, error)
}
type userDbImpl struct {
	DB         *mongo.Database
	Collection *mongo.Collection
}

func NewUserDb(DB *mongo.Database, collection *mongo.Collection) UserDb {
	return &userDbImpl{DB: DB, Collection: collection}
}
func (u *userDbImpl) Add(ctx context.Context, user domain.User) (string, error) {
	resp, err := u.Collection.InsertOne(ctx, user)
	return resp.InsertedID.(string), err
}
func (u *userDbImpl) Deactive(ctx context.Context, user domain.User) error {
	one, err := u.Collection.UpdateOne(ctx, bson.M{"_id": user.Email}, bson.D{{"$set", bson.D{{"is_active", false}}}})
}
