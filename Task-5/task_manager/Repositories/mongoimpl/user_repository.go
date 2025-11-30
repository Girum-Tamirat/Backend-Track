package mongoimpl

import (
	"context"
	"errors"

	"task_manager/Domain"
	"task_manager/Repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepo struct {
	coll *mongo.Collection
}

func NewUserRepository(client *MongoClient) Repositories.UserRepository {
	coll := client.Client.Database(client.DBName).Collection("users")
	_, _ = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return &userRepo{coll: coll}
}

func (r *userRepo) Create(ctx context.Context, u Domain.User) (Domain.User, error) {
	_, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		return Domain.User{}, err
	}
	u.Password = ""
	return u, nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (Domain.User, error) {
	var u Domain.User
	if err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return Domain.User{}, errors.New("not found")
		}
		return Domain.User{}, err
	}
	return u, nil
}

func (r *userRepo) UpdateRole(ctx context.Context, username, role string) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"role": role}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *userRepo) CountUsers(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}
