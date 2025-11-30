package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"-"` // hashed
	Role     string             `bson:"role" json:"role"` // "admin" or "user"
}
