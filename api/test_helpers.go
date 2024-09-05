package api

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/assiljaby/gotel-reservation/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	DBNAME := os.Getenv("MONGO_DB_NAME")
	if err := tdb.client.Database(DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	dburi := os.Getenv("MONGO_DB_URL_TEST")
	DBNAME := os.Getenv("MONGO_DB_NAME")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, DBNAME)
	return &testdb{
		client: client,
		Store: &db.Store{
			Hotel:   hotelStore,
			User:    db.NewMongoUserStore(client, DBNAME),
			Room:    db.NewMongoRoomStore(client, hotelStore, DBNAME),
			Booking: db.NewMongoBookingStore(client, DBNAME),
		},
	}
}
