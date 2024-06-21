package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dclouisDan/chat-app-api/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	app := fiber.New()

	// Register routes with Fiber app
	handler.RegisterRoutes(app)

	t.Run("should fail if user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName:  "client",
			Email:     "invalid",
			Password:  "asd",
		}
		marshalled, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		req.Header.Set("Content-Type", "application/json") // Set content type JSON

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

    assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should correctly register the user", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName:  "client",
			Email:     "user@email.com",
			Password:  "asd",
		}
		marshalled, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		req.Header.Set("Content-Type", "application/json") // Set content type JSON

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
    
    assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	})

  t.Run("should fail if invalid credentials are provided", func(t *testing.T) {
	payload := types.LoginUserPayload{
		Email:    "user@email.com",
		Password: "wrongpassword",
	}
	marshalled, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err, "error testing request")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
})

}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(user types.User) error {
	return nil
}

