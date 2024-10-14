package models

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	Activated     bool      `json:"activated"`
	Bio           string    `json:"bio,omitempty"`
	Password      password  `json:"-"`
	ProfilePicUrl string    `json:"profilePicUrl,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type password struct {
	hash      []byte
	plainText string
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Create(user User) error {

	q := `
	INSERT INTO users (first_name, last_name, email, password_hash)
	VALUES ($1, $2, $3, $4)
	`
	args := []any{user.FirstName, user.LastName, user.Email, user.Password.hash}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, args...).Scan(args...)
	if err != nil {
		return err
	}
	return nil

}

//func (m UserModel) Create(user *User) error {
//	q := `
//	INSERT INTO users
//	VALUES ($1, $2, $3, $4, $5, $6, $7)
//
//`
//}
