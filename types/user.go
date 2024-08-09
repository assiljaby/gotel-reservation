package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
	minFirstNameLen = 3
	minLastNameLen = 3
	minPasswordLen = 8
)

type UserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email  	  string `json:"email"`
	Password  string `json:"password"`
}

func (u UserParams) Validate() error {
	if fl := len(u.FirstName); fl < minFirstNameLen {
		return fmt.Errorf("expected first name to be at least %d characters long, got %d instead", minFirstNameLen, fl)
	}
	if ll := len(u.LastName); ll < minLastNameLen {
		return fmt.Errorf("expected last name to be at least %d characters long, got %d instead", minLastNameLen, ll)
	}
	if pl := len(u.Password); pl < minPasswordLen {
		return fmt.Errorf("expected password to be at least %d characters long, got %d instead", minPasswordLen, pl)
	}
	if !isEmailValid(u.Email) {
		return fmt.Errorf("invalid email: %s", u.Email)
	}
	return nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(e)
} 

type User struct {
	ID 		  	 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName 	 string `bson:"firstName" json:"firstName"`
	LastName  	 string `bson:"lastName" json:"lastName"`
	Email  	  	 string `bson:"email" json:"email"`
	PasswordHash string `bson:"passwordHash" json:"passwordHash"`
}

func NewUserFromParams(userPrms UserParams) (*User, error) {
	err := userPrms.Validate()
	if err != nil {
		return nil, err
	}
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