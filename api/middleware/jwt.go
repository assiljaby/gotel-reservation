package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/assiljaby/gotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return fmt.Errorf("unauthorized")
		}

		claims, err := validateToken(token[0])
		if err != nil {
			return fmt.Errorf("unauthorized")
		}

		validUntillFloat := claims["validUntill"].(float64)
		validUntill := int64(validUntillFloat)
		if time.Now().Unix() > validUntill {
			return fmt.Errorf("token expired")
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserById(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}
func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method:", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil
}
