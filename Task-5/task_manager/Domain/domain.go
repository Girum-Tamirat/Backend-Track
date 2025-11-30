package Domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Domain entities — independent of frameworks

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password,omitempty" json:"-"`
	Role     string             `bson:"role" json:"role"` // "admin" or "user"
}

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	DueDate     string             `bson:"due_date,omitempty" json:"due_date,omitempty"`
	Status      string             `bson:"status" json:"status"`
}
