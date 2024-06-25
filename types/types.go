package types

import (
	"database/sql"
	"time"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
  UpdateUser(User) error
  UpdateUserProfilePicture(userID int, path string) error
}

type User struct {
	ID             int            `json:"id"`
	FirstName      string         `json:"firstName"`
	LastName       string         `json:"lastName"`
	Email          string         `json:"email"`
	Password       string         `json:"password"`
	ProfilePicture sql.NullString `json:"profilePicture"`
	CreatedAt      time.Time      `json:"createdAt"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type UpdateUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
