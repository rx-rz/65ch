package data

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}
type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m UserModel) Create(user User) error {
	q := `
	INSERT INTO users (first_name, last_name, email, password_hash)
	VALUES ($1, $2, $3, $4)
	RETURNING first_name, last_name , email, created_at
	`
	args := []any{&user.FirstName, &user.LastName, &user.Email, &user.Password}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, args...).Scan(args...)
	return determineDBError(err)
}

func (m UserModel) FindByEmail(email string) (*User, error) {
	var user User
	q := `
	SELECT first_name, last_name, email, id, password_hash FROM users WHERE email = $1
	`
	args := []any{&user.FirstName, &user.LastName, &user.Email, &user.ID, &user.Password}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, email).Scan(args...)
	return &user, determineDBError(err)
}

//func (m UserModel) Create(user *User) error {
//	q := `
//	INSERT INTO users
//	VALUES ($1, $2, $3, $4, $5, $6, $7)
//
//`
//}
