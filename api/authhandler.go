package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/assiljaby/gotel-reservation/db"
	"github.com/assiljaby/gotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(store *db.Store) *AuthHandler {
	return &AuthHandler{
		userStore: store.User,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResponse struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResponse{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	if !types.IsPasswordValid(user.PasswordHash, params.Password) {
		return invalidCredentials(c)
	}

	token := CreateTokenFromUser(user)
	if len(token) == 0 {
		return fmt.Errorf("could not auth")
	}

	res := AuthResponse{
		User:  user,
		Token: token,
	}

	fmt.Println("Authed user:", user.Email)
	return c.JSON(res)
}

func CreateTokenFromUser(user *types.User) string {
	validUntill := time.Now().Add(time.Hour * 4).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          user.ID,
		"email":       user.Email,
		"validUntill": validUntill,
	})

	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("could not sign token")
	}
	return tokenStr
}
