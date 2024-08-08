package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/assiljaby/gotel-reservation/api"
	"github.com/assiljaby/gotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "gotel-reservation"
const userCollection = "users"


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongodbURI := os.Getenv("MONGODB_URI")
	// mongodbAtlasURI := os.Getenv("MONGODB_ATLAS_URI")
	if mongodbURI == "" {
		log.Fatal("Database URI was not set correctly.")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal("Failed to connect to database.")
	}

	// ctx := context.Background()
	coll := client.Database(dbName).Collection(userCollection)

	scrubloard := types.User{
		FirstName: "Scrub",
		LastName: "Lord",
	}

	res, err := coll.InsertOne(context.TODO(), scrubloard)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	listenPort := flag.String("listenPort", ":3000", "The server is listening to this port")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	app.Get("/foo", handleFoo)
	apiv1.Get("/users", api.HandleGetUsers)
	apiv1.Get("/users/:id", api.HandleGetUser)

	app.Listen(*listenPort)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Bar!"})
}