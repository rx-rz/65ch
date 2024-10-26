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
	ID            string    `db:"id"`
	Email         string    `db:"email"`
	Password      string    `db:"password"`
	FirstName     string    `db:"first_name"`
	Bio           string    `db:"bio"`
	ProfilePicUrl string    `db:"profile_picture_url"`
	ResetToken    *string   `db:"reset_token"`
	LastName      string    `db:"last_name"`
	Activated     bool      `db:"activated"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (m UserModel) Create(user *User) error {
	q := `
	INSERT INTO users (first_name, last_name, email, password_hash, bio, profile_picture_url)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING first_name, last_name , email, bio, profile_picture_url, created_at
	`
	if user.Bio == "" {
		user.Bio = "Enter your bio"
	}
	if user.ProfilePicUrl == "" {
		user.ProfilePicUrl = "https://placehold.co/400?text=U"
	}
	args := []any{user.FirstName, user.LastName, user.Email, user.Password, user.Bio, user.ProfilePicUrl}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, args...).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Bio, &user.ProfilePicUrl, &user.CreatedAt)
	return DetermineDBError(err, "user_create")
}

func (m UserModel) FindByEmail(email string) (*User, error) {
	var user User
	q := `
	SELECT first_name, last_name, email, id, password_hash, bio, profile_picture_url FROM users WHERE email = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, email).Scan(&user.FirstName, &user.LastName, &user.Email, &user.ID, &user.Password, &user.Bio, &user.ProfilePicUrl)
	if err != nil {
		return nil, DetermineDBError(err, "user_findbyemail")
	}
	return &user, nil
}

func (m UserModel) FindByID(id string) (*User, error) {
	var user User
	q := `
	SELECT first_name, last_name, email, id, password_hash, bio, profile_picture_url FROM users WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, id).Scan(&user.FirstName, &user.LastName, &user.Email, &user.ID, &user.Password, &user.Bio, &user.ProfilePicUrl)
	if err != nil {
		return nil, DetermineDBError(err, "user_findbyid")
	}
	return &user, nil
}

func (m UserModel) UpdateDetails(user *User) (*User, error) {
	var userDetails User
	q := `
	UPDATE users SET first_name = $1, last_name = $2, bio = $3, profile_picture_url = $4, activated = $5
	WHERE id = $6
	RETURNING id, first_name, last_name, bio, profile_picture_url, activated
	`
	args := []any{&user.FirstName, &user.LastName, &user.Bio, &user.ProfilePicUrl, &user.Activated, &user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, args...).Scan(&userDetails.ID, &userDetails.FirstName, &userDetails.LastName, &userDetails.Bio, &userDetails.ProfilePicUrl, &userDetails.Activated)
	if err != nil {
		return nil, DetermineDBError(err, "user_updatedetails")
	}
	return &userDetails, nil
}

func (m UserModel) UpdateEmail(email, newEmail string) error {
	q := `
	UPDATE users SET email = $1 where email = $2
`
	args := []any{newEmail, email}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, q, args...)
	return DetermineDBError(err, "user_updateemail")
}

func (m UserModel) UpdatePassword(email, newPassword string) error {
	q := `
	UPDATE users SET password_hash = $1 where email = $2
	`
	args := []any{newPassword, email}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, q, args...)
	return DetermineDBError(err, "user_updatepassword")
}
