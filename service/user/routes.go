package user

import (
	"fmt"
	"log"
	"time"

	"github.com/dclouisDan/chat-app-api/config"
	"github.com/dclouisDan/chat-app-api/service/auth"
	"github.com/dclouisDan/chat-app-api/types"
	"github.com/dclouisDan/chat-app-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/login", h.handleLogin)
	router.Post("/register", h.handleRegister)
	router.Get("/profile", auth.WithJWTAuth(h.handleProfile, h.store))
}

// User Login
func (h *Handler) handleLogin(c *fiber.Ctx) error {
	var payload types.LoginUserPayload

	// parse payload
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Parse error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// validate payload
	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("invalid payload %v", errors),
		})
	}

	// get user by email
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "not found, invalid email or password",
		})
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "not found, invalid email or password",
		})
	}

	secret := []byte(config.Envs.JWTSecret)
	token, expireAt, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Unix(expireAt, 0),
		HTTPOnly: true,
	})
  
  return c.Status(fiber.StatusOK).JSON(fiber.Map{
    "token": token,
  })
}

// User Registration
func (h *Handler) handleRegister(c *fiber.Ctx) error {
	var payload types.RegisterUserPayload

	// parse payload
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Parse error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// validate payload
	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error" : fmt.Sprintf("invalid payload %v", errors),
    })
	}

	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error" : fmt.Sprintf("user with email %s already exists", payload.Email),
    })
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error" : err.Error(),
    })
	}

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error" : err.Error(),
    })
	}

  
  return c.Status(fiber.StatusCreated).JSON(fiber.Map{
    "message" : "user created.",
  })
}

// Account profile
func (h *Handler) handleProfile(c *fiber.Ctx) error {
	userID := auth.GetIDFromContext(c)

	u, err := h.store.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error" : "user not found",
    })
	}
  
  return c.Status(fiber.StatusOK).JSON(u)
}
