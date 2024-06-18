package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userID int) (string, error) {
  expiration := time.Second * time.Duration(1000)

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "userID" : strconv.Itoa(userID),
    "expiredAt" : time.Now().Add(expiration).Unix(),
  })

  tokenString, err := token.SignedString(secret)
  if err != nil {
    return "", nil

  }
    return tokenString, nil
}
