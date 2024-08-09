package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
)

type UserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email  	  string `json:"email"`
	Password  string `json:"password"`
}
type User struct {
	ID 		  	 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName 	 string `bson:"firstName" json:"firstName"`
	LastName  	 string `bson:"lastName" json:"lastName"`
	Email  	  	 string `bson:"email" json:"email"`
	PasswordHash string `bson:"passwordHash" json:"passwordHash"`
}

func NewUserFromParams(userPrms UserParams) (*User, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(userPrms.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: userPrms.FirstName,
		LastName: userPrms.LastName,
		Email: userPrms.Email,
		PasswordHash: string(passHash),
	}, nil
}