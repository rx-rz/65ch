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

func (m *UserModel) Create(user *User) (*User, error) {
	const query = `
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	newUser := &User{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Bio,
		user.ProfilePicUrl,
	).Scan(
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.Email,
		&newUser.Bio,
		&newUser.ProfilePicUrl,
		&newUser.CreatedAt,
	)
	return newUser, DetermineDBError(err, "user_create")
}

func (m *UserModel) FindByEmail(email string) (*User, error) {
	const query = `
	SELECT first_name, last_name, email, id, password_hash, bio, profile_picture_url
	FROM users
	WHERE email = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	user := &User{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.ID,
		&user.Password,
		&user.Bio,
		&user.ProfilePicUrl,
	)
	if err != nil {
		return nil, DetermineDBError(err, "user_findbyemail")
	}
	return user, nil
}

func (m *UserModel) FindByID(id string) (*User, error) {
	const query = `
	SELECT first_name, last_name, email, id, password_hash, bio, profile_picture_url
	FROM users
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	user := &User{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.ID,
		&user.Password,
		&user.Bio,
		&user.ProfilePicUrl,
	)
	if err != nil {
		return nil, DetermineDBError(err, "user_findbyemail")
	}
	return user, nil
}

func (m *UserModel) UpdateDetails(user *User) (*User, error) {
	const query = `
	UPDATE users SET 
		first_name = $1, 
		last_name = $2, 
		bio = $3, 
		profile_picture_url = $4, 
		activated = $5
	WHERE id = $6
	RETURNING id, first_name, last_name, bio, profile_picture_url, activated
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	updatedUser := &User{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Bio,
		user.ProfilePicUrl,
		user.Activated,
		user.ID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Bio,
		&updatedUser.ProfilePicUrl,
		&updatedUser.Activated)

	if err != nil {
		return nil, DetermineDBError(err, "user_updatedetails")
	}
	return updatedUser, nil
}

func (m *UserModel) UpdateEmail(email, newEmail string) (*ModifiedData, error) {
	const query = `
	UPDATE users 
	SET email = $1, updated_at = $2
	WHERE email = $3
	RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	updatedUser := &User{}
	updateTimestamp := time.Now().UTC()

	err := m.DB.QueryRowContext(
		ctx,
		query,
		newEmail,
		updateTimestamp,
		email,
	).Scan(
		&updatedUser.ID,
	)
	if err != nil {
		return nil, DetermineDBError(err, "user_updateemail")
	}

	return &ModifiedData{
			ID:        updatedUser.ID,
			Timestamp: updateTimestamp,
		},
		DetermineDBError(err, "user_updateemail")
}

func (m *UserModel) UpdatePassword(email, newPassword string) (*ModifiedData, error) {
	q := `
	UPDATE users 
	SET password_hash = $1, updated_at = $2
	WHERE email = $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	updatedUser := &User{}
	updateTimestamp := time.Now().UTC()

	err := m.DB.QueryRowContext(
		ctx,
		q,
		newPassword,
		updateTimestamp,
		email,
	).Scan(
		&updatedUser.ID,
	)

	if err != nil {
		return nil, DetermineDBError(err, "user_updatepassword")
	}

	return &ModifiedData{
			ID:        updatedUser.ID,
			Timestamp: updateTimestamp,
		},
		DetermineDBError(err, "user_updatepassword")
}
