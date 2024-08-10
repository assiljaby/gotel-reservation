package db

import (
	"context"
	"fmt"

	"github.com/assiljaby/gotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type UserStore interface {
	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	CreateUser(context.Context, *types.UserWithoutID) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, types.UserWithoutID) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll: client.Database(DBNAME).Collection(userColl),
	}
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	var user types.User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.UserWithoutID) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser := types.User{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		PasswordHash: user.PasswordHash,
	}

	createdUser.ID = res.InsertedID.(primitive.ObjectID)
	return &createdUser, nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, user types.UserWithoutID) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	
	return nil
}


func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("no users were deleted by this operation")
	}
	return nil
}