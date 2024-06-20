package auth

import "testing"

func TestCreateJWT(t *testing.T) {
  secret := []byte("secret")

  token, _, err := CreateJWT(secret, 1)
  if err != nil {
    t.Errorf("error creating JWT: %v", err)
  }

  if token == "" {
    t.Error("expected token to be not empty")
  }
}
