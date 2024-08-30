package main

import (
	"context"
	"log"
	"os"

	"github.com/assiljaby/gotel-reservation/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Loading Env Vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Geting MongoDB URI from ENV
	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		log.Fatal("Database URI was not set correctly.")
	}
	DBNAME := os.Getenv("MONGO_DB_NAME")
	if mongodbURI == "" {
		log.Fatal("Database URI was not set correctly.")
	}

	// Initializing DB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal("Failed to connect to database.")
	}

	hotelStore := db.NewMongoHotelStore(client, DBNAME)

	_ = hotelStore
}
