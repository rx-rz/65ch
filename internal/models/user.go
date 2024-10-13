package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Activated     bool      `json:"activated"`
	Bio           string    `json:"bio,omitempty"`
	Password      string    `json:"password"`
	ProfilePicUrl string    `json:"profilePicUrl,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type UserModel struct {
	DB *sql.DB
}

//func (m UserModel) Create(user *User) error {
//	q := `
//	INSERT INTO users
//	VALUES ($1, $2, $3, $4, $5, $6, $7)
//
//`
//}
