package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task stored in MongoDB
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" binding:"required"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	DueDate     string             `bson:"due_date,omitempty" json:"due_date,omitempty"`
	Status      string             `bson:"status" json:"status" binding:"required"`
}
