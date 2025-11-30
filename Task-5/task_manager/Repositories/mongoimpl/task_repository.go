package mongoimpl

import (
	"context"
	"errors"

	"task_manager/Domain"
	"task_manager/Repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type taskRepo struct {
	coll *mongo.Collection
}

func NewTaskRepository(client *MongoClient) Repositories.TaskRepository {
	coll := client.Client.Database(client.DBName).Collection("tasks")
	return &taskRepo{coll: coll}
}

func (r *taskRepo) Create(ctx context.Context, t Domain.Task) (Domain.Task, error) {
	res, err := r.coll.InsertOne(ctx, t)
	if err != nil {
		return Domain.Task{}, err
	}
	var created Domain.Task
	err = r.coll.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&created)
	return created, err
}

func (r *taskRepo) FindAll(ctx context.Context) ([]Domain.Task, error) {
	cur, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []Domain.Task
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *taskRepo) FindByID(ctx context.Context, id primitive.ObjectID) (Domain.Task, error) {
	var t Domain.Task
	if err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return Domain.Task{}, errors.New("not found")
		}
		return Domain.Task{}, err
	}
	return t, nil
}

func (r *taskRepo) Update(ctx context.Context, id primitive.ObjectID, patch map[string]interface{}) (Domain.Task, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": patch}, opts)
	var updated Domain.Task
	if err := res.Decode(&updated); err != nil {
		return Domain.Task{}, err
	}
	return updated, nil
}

func (r *taskRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("not found")
	}
	return nil
}
