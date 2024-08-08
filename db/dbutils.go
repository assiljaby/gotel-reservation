package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const DBNAME = "gotel-reservation"

func ParseObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	return oid
}