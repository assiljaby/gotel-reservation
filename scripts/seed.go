package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/assiljaby/gotel-reservation/api"
	"github.com/assiljaby/gotel-reservation/db"
	"github.com/assiljaby/gotel-reservation/db/fixtures"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/rand"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		ctx           = context.Background()
		mongoEndpoint = os.Getenv("MONGODB_URI")
		DBNAME        = os.Getenv("MONGO_DB_NAME")
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, DBNAME)
	store := &db.Store{
		User:    db.NewMongoUserStore(client, DBNAME),
		Booking: db.NewMongoBookingStore(client, DBNAME),
		Room:    db.NewMongoRoomStore(client, hotelStore, DBNAME),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Println("james ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "some hotel", "bermuda", 5, nil)
	room := fixtures.AddRoom(store, "large", true, 88.44, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ->", booking.ID)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}
