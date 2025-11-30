package data

import (
	"context"
	"errors"
	"time"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	coll   *mongo.Collection
	client *mongo.Client
	timeout time.Duration
}

var userSvc *UserService

// InitUserService initializes user service (call once)
func InitUserService(uri, dbName string) error {
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
	coll := client.Database(dbName).Collection("users")
	// unique index on username
	_, _ = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	userSvc = &UserService{coll: coll, client: client, timeout: 5 * time.Second}
	return nil
}

func GetUserService() *UserService {
	return userSvc
}

func (s *UserService) Close() error {
	if s.client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}

// CreateUser registers a new user. If DB has no users, first user becomes admin.
func (s *UserService) CreateUser(username, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	// check uniqueness and count
	count, err := s.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	role := "user"
	if count == 0 {
		role = "admin" // first user -> admin
	}
	user := models.User{Username: username, Password: string(hashed), Role: role}
	_, err = s.coll.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.New("username already exists")
		}
		return nil, err
	}
	// hide password in returned user
	user.Password = ""
	return &user, nil
}

// AuthenticateUser verifies username/password and returns user (without password) or error
func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var user models.User
	if err := s.coll.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	user.Password = ""
	return &user, nil
}

// PromoteUser sets role to admin. Only admins call this.
func (s *UserService) PromoteUser(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
    
	res, err := s.coll.UpdateOne(
        ctx, 
        bson.M{"username": username}, 
        bson.M{"$set": bson.M{"role": "admin"}}, // <-- $set operator added here
    )
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

// GetByUsername returns user with role (no password)
func (s *UserService) GetByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	var user models.User
	if err := s.coll.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return nil, err
	}
	user.Password = ""
	return &user, nil
}