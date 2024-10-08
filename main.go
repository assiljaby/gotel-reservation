package main

import (
	"context"
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
	ErrorHandler: api.ErrorHandler,
}

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

	// Handler Initialization
	var (
		hotelStore   = db.NewMongoHotelStore(client, DBNAME)
		userStore    = db.NewMongoUserStore(client, DBNAME)
		roomStore    = db.NewMongoRoomStore(client, hotelStore, DBNAME)
		bookingStore = db.NewMongoBookingStore(client, DBNAME)
		store        = &db.Store{
			Hotel:   hotelStore,
			User:    userStore,
			Room:    roomStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		AuthHandler    = api.NewAuthHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", api.JWTAuth(userStore))
		admin          = apiv1.Group("/admin", api.AdminAuth)
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
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// Room Handlers
	apiv1.Post("/rooms/:id/book", roomHandler.HandleBookRoom)
	apiv1.Post("/rooms", roomHandler.HandleGetRooms)

	// Booking Handlers
	apiv1.Get("/bookings/:id", bookingHandler.HandleGetBooking)
	apiv1.Put("/bookings/:id/cancel", bookingHandler.HandleCancelBooking)

	// Admin routes
	admin.Get("/bookings", bookingHandler.HandleGetBookings)

	listenPort := os.Getenv("LISTEN_PORT")
	app.Listen(listenPort)
}
