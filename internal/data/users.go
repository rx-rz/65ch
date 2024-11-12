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

func (m *UserModel) Create(ctx context.Context, user *User) (*User, error) {
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
	if err != nil {
		return nil, DetermineDBError(err, "user_create")
	}
	return newUser, nil
}

func (m *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	const query = `
	SELECT first_name, last_name, email, id, password_hash, bio, profile_picture_url
	FROM users
	WHERE email = $1
	`
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

func (m *UserModel) GetByID(ctx context.Context, id string) (*User, error) {
	const query = `
	SELECT first_name, last_name, email, id, password_hash, bio, profile_picture_url
	FROM users
	WHERE id = $1
	`

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

func (m *UserModel) UpdateDetails(ctx context.Context, user *User) (*User, error) {
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

func (m *UserModel) UpdateEmail(ctx context.Context, email, newEmail string) (*ModifiedData, error) {
	const query = `
	UPDATE users 
	SET email = $1, updated_at = $2
	WHERE email = $3
	RETURNING id
	`

	data := &ModifiedData{}
	updateTimestamp := time.Now().UTC()

	err := m.DB.QueryRowContext(
		ctx,
		query,
		newEmail,
		updateTimestamp,
		email,
	).Scan(
		&data.ID,
	)
	if err != nil {
		return nil, DetermineDBError(err, "user_updateemail")
	}
	data.Timestamp = updateTimestamp
	return data, DetermineDBError(err, "user_updateemail")
}

func (m *UserModel) UpdatePassword(ctx context.Context, email, newPassword string) (*ModifiedData, error) {
	q := `
	UPDATE users 
	SET password_hash = $1, updated_at = $2
	WHERE email = $3
	`

	data := &ModifiedData{}
	updateTimestamp := time.Now().UTC()

	err := m.DB.QueryRowContext(
		ctx,
		q,
		newPassword,
		updateTimestamp,
		email,
	).Scan(
		&data.ID,
	)

	if err != nil {
		return nil, DetermineDBError(err, "user_updatepassword")
	}
	data.Timestamp = updateTimestamp
	return data, DetermineDBError(err, "user_updatepassword")
}
