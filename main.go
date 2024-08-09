package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/assiljaby/gotel-reservation/api"
	"github.com/assiljaby/gotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	// Geetting and parsing flags
	listenPort := flag.String("listenPort", ":3000", "The server is listening to this port")
	flag.Parse()

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

	// Initializing DB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal("Failed to connect to database.")
	}

	// Handler Initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)

	apiv1.Post("/users", userHandler.HandlePostUser)

	app.Listen(*listenPort)
}