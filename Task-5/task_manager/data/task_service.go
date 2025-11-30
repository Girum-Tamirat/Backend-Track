package data

import (
	"context"
	"errors"
	"time"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskService struct {
	Collection *mongo.Collection
}

func NewTaskService() (*TaskService, error) {

	// Note: For production, use environment variables for URI and handle connection pooling
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	
	// Ping to confirm connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.New("failed to ping MongoDB: " + err.Error())
	}

	db := client.Database("taskdb")
	collection := db.Collection("tasks")

	return &TaskService{Collection: collection}, nil
}

func (s *TaskService) Create(task models.Task) (*mongo.InsertOneResult, error) {
	ctx := context.Background()
	return s.Collection.InsertOne(ctx, task)
}

func (s *TaskService) GetAll() ([]models.Task, error) {
	ctx := context.Background()
	cursor, err := s.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) GetByID(id string) (models.Task, error) {
	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, errors.New("invalid task ID format")
	}

	var task models.Task
	err = s.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	if err == mongo.ErrNoDocuments {
		return models.Task{}, errors.New("task not found")
	}
	return task, err
}

func (s *TaskService) Update(id string, data models.Task) (models.Task, error) {
	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, errors.New("invalid task ID format")
	}

	// This is the crucial part: using the $set operator
	update := bson.M{
		"$set": bson.M{
			"title":       data.Title,
			"description": data.Description,
			"due_date":    data.DueDate,
			"status":      data.Status,
		},
	}

	res, err := s.Collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return models.Task{}, err
	}
	if res.MatchedCount == 0 {
		return models.Task{}, errors.New("task not found")
	}

	// Retrieve the updated document for response
	var updatedTask models.Task
	err = s.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedTask)
	return updatedTask, err
}

func (s *TaskService) Delete(id string) error {
	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID format")
	}

	res, err := s.Collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}