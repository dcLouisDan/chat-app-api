package auth

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dclouisDan/chat-app-api/config"
	"github.com/dclouisDan/chat-app-api/types"
	"github.com/dclouisDan/chat-app-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userID int) (string, int64, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	expireAt := time.Now().Add(expiration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": expireAt,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", 0, nil
	}
	return tokenString, expireAt, nil
}

func WithJWTAuth(handlerFunc fiber.Handler, store types.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := utils.GetTokenFromRequest(c)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %s", err.Error())
			return permissionDenied(c)
		}

		if !token.Valid {
			log.Println("invalid token")
			return permissionDenied(c)
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims[string(UserKey)].(string)

		userID, _ := strconv.Atoi(str)
		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			return permissionDenied(c)
		}

		ctx := context.WithValue(c.UserContext(), UserKey, u.ID)
		c.SetUserContext(ctx)

		return handlerFunc(c)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Permission Denied",
	})
}

func GetIDFromContext(c *fiber.Ctx) int {
	userID, ok := c.UserContext().Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
