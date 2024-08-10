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

func (u UserParams) Validate() map[string]string {
	errors := map[string]string{}
	if fl := len(u.FirstName); fl < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("expected first name to be at least %d characters long, got %d instead", minFirstNameLen, fl)
	}
	if ll := len(u.LastName); ll < minLastNameLen {
		errors["lastName"] =  fmt.Sprintf("expected last name to be at least %d characters long, got %d instead", minLastNameLen, ll)
	}
	if pl := len(u.Password); pl < minPasswordLen {
		errors["password"] =  fmt.Sprintf("expected password to be at least %d characters long, got %d instead", minPasswordLen, pl)
	}
	if !isEmailValid(u.Email) {
		errors["email"] =  fmt.Sprintf("invalid email: %s", u.Email)
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(e)
} 

type UserWithoutID struct {
	FirstName 	 string `bson:"firstName" json:"firstName"`
	LastName  	 string `bson:"lastName" json:"lastName"`
	Email  	  	 string `bson:"email" json:"email"`
	PasswordHash string `bson:"passwordHash" json:"passwordHash"`
}

type User struct {
	ID 		  	 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserWithoutID
}

func NewUserFromParams(userPrms UserParams) (*UserWithoutID, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(userPrms.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &UserWithoutID{
		FirstName: userPrms.FirstName,
		LastName: userPrms.LastName,
		Email: userPrms.Email,
		PasswordHash: string(passHash),
	}, nil
}

