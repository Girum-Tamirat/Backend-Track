package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Title       string             `json:"title" bson:"title" binding:"required"`
    Description string             `json:"description" bson:"description,omitempty"`
    DueDate     string             `json:"due_date" bson:"due_date,omitempty"`
    Status      string             `json:"status" bson:"status" binding:"required"`
}
