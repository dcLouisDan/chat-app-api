package user

import (
	"database/sql"
	"fmt"

	"github.com/dclouisDan/chat-app-api/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("User not found.")
	}

	return u, nil
}

func (s *Store) GetUserByID(userID int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", userID)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("User not found.")
	}

	return u, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateUserProfilePicture(userID int, path string) error {
	_, err := s.db.Exec("UPDATE users SET profilePicture = ? WHERE id = ?;", path, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateUser(user types.User) error {
	_, err := s.db.Exec("UPDATE users SET firstname = ?, lastName = ?, email = ? WHERE id = ?;", user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.ProfilePicture,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
