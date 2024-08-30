package db

import (
	"context"

	"github.com/assiljaby/gotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	CreateHotel(context.Context, types.Hotel) (*types.HotelWithoutID, error)
	Update(context.Context, Map, Map) error
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbName string) *MongoHotelStore {
	return &MongoHotelStore{
		client:     client,
		collection: client.Database(dbName).Collection("hotels"),
	}
}

func (s *MongoHotelStore) CreateHotel(ctx context.Context, hotel types.HotelWithoutID) (*types.Hotel, error) {
	res, err := s.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	return &types.Hotel{
		ID:       res.InsertedID.(primitive.ObjectID),
		Name:     hotel.Name,
		Location: hotel.Location,
		Rooms:    hotel.Rooms,
		Rating:   hotel.Rating,
	}, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filter Map, update Map) error {
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}
