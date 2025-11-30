package data

import (
	"context"
	"errors"
	"time"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskService struct {
	client *mongo.Client
	coll   *mongo.Collection
	timeout time.Duration
}

var taskSvc *TaskService

func InitTaskService(uri, dbName string) error {
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	if dbName == "" {
		dbName = "taskdb"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}
	coll := client.Database(dbName).Collection("tasks")
	taskSvc = &TaskService{client: client, coll: coll, timeout: 5 * time.Second}
	// index for status/title if needed (optional)
	_, _ = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bson.D{{Key: "title", Value: 1}}})
	return nil
}

func GetTaskService() *TaskService {
	return taskSvc
}

func (s *TaskService) Close() error {
	if s.client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}

func (s *TaskService) CreateTask(t models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	res, err := s.coll.InsertOne(ctx, t)
	if err != nil {
		return models.Task{}, err
	}
	// fetch inserted document to return with ID
	var inserted models.Task
	if err := s.coll.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&inserted); err != nil {
		return models.Task{}, err
	}
	return inserted, nil
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) GetTaskByID(filter bson.M) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	var task models.Task
	if err := s.coll.FindOne(ctx, filter).Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, err
	}
	return task, nil
}

// UpdateTask takes the filter (ID) and the fields to update (update)
func (s *TaskService) UpdateTask(filter bson.M, update bson.M) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	updateDoc := bson.M{"$set": update}

	// Use FindOneAndUpdate with $set to update fields and return the new document
	res := s.coll.FindOneAndUpdate(ctx, filter, updateDoc, options.FindOneAndUpdate().SetReturnDocument(options.After))
	var updated models.Task
	if err := res.Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, err
	}
	return updated, nil
}

func (s *TaskService) DeleteTask(filter bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	res, err := s.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}