package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/assiljaby/gotel-reservation/api"
	"github.com/assiljaby/gotel-reservation/api/middleware"
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

	DBNAME := os.Getenv("MONGO_DB_NAME")
	if mongodbURI == "" {
		log.Fatal("Database URI was not set correctly.")
	}

	// Initializing DB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal("Failed to connect to database.")
	}

	// Handler Initialization
	var (
		hotelStore = db.NewMongoHotelStore(client, DBNAME)
		userStore  = db.NewMongoUserStore(client, DBNAME)
		roomStore  = db.NewMongoRoomStore(client, hotelStore, DBNAME)
		store      = &db.Store{
			Hotel: hotelStore,
			User:  userStore,
			Room:  roomStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		RoomHandler  = api.NewRoomHandler(store)
		AuthHandler  = api.NewAuthHandler(store)
		app          = fiber.New(config)
		auth         = app.Group("/api")
		apiv1        = app.Group("/api/v1", middleware.JWTAuth(userStore))
	)

	// Auth Handlers
	auth.Post("/auth", AuthHandler.HandleAuthenticate)
	auth.Post("/users", userHandler.HandlePostUser)

	// User Handlers
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Put("/users/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)

	// Hotel Handlers
	apiv1.Get("hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// Booking Handlers
	apiv1.Post("room/:id/book", RoomHandler.HandleBookRoom)

	app.Listen(*listenPort)
}
