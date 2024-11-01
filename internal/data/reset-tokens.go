package data

import (
	"context"
	"database/sql"
	"time"
)

type ResetTokenModel struct {
	DB *sql.DB
}

type ResetToken struct {
	ID         int       `json:"id"`
	UserID     string    `json:"user_id"`
	ResetToken string    `json:"reset_token"`
	Expiration time.Time `json:"expiration"`
}

func (m ResetTokenModel) Create(ctx context.Context, resetToken *ResetToken) (*ResetToken, error) {
	const query = `
	INSERT INTO  reset_tokens (user_id, reset_token, expiration) 
	VALUES ($1, $2, $3)
	RETURNING id, user_id, reset_token, expiration
	`
	newResetToken := &ResetToken{}

	err := m.DB.QueryRowContext(
		ctx,
		query,
		resetToken.UserID,
		resetToken.ResetToken,
		resetToken.Expiration,
	).Scan(
		&newResetToken.ID,
		&newResetToken.UserID,
		&newResetToken.ResetToken,
		&newResetToken.Expiration,
	)
	if err != nil {
		return nil, DetermineDBError(err, "resettoken_create")
	}
	return newResetToken, nil

}

func (m ResetTokenModel) GetByUserID(ctx context.Context, userId string) (*ResetToken, error) {
	const query = `
	SELECT id, user_id, reset_token, expiration 
	FROM reset_tokens 
	WHERE user_id = $1`

	resetToken := &ResetToken{}

	err := m.DB.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&resetToken.ID,
		&resetToken.UserID,
		&resetToken.ResetToken,
		&resetToken.Expiration,
	)
	if err != nil {
		return nil, DetermineDBError(err, "resettoken_getbyuserid")
	}
	return resetToken, nil
}

func (m ResetTokenModel) GetByToken(ctx context.Context, token string) (*ResetToken, error) {
	const query = `
	SELECT id, user_id, reset_token, expiration 
	FROM reset_tokens 
	WHERE reset_token = $1
	`

	resetToken := &ResetToken{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		token,
	).Scan(
		&resetToken.ID,
		&resetToken.UserID,
		&resetToken.ResetToken,
		&resetToken.Expiration,
	)

	if err != nil {
		return nil, DetermineDBError(err, "resettoken_getbytoken")
	}
	return resetToken, nil
}

func (m ResetTokenModel) Update(ctx context.Context, resetToken *ResetToken) (*ModifiedData, error) {
	const query = `
	UPDATE reset_tokens 
	SET reset_token = $1, expiration = $2 
	WHERE user_id = $3
	RETURNING id
	`
	data := &ModifiedData{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		resetToken.ResetToken,
		resetToken.Expiration,
		resetToken.UserID,
	).Scan(
		&data.ID,
	)
	data.Timestamp = time.Now().UTC()
	if err != nil {
		return nil, DetermineDBError(err, "resettoken_update")
	}
	return data, nil
}

func (m ResetTokenModel) DeleteByUserId(ctx context.Context, userId string) (*ModifiedData, error) {
	const query = `
	DELETE FROM reset_tokens 
	WHERE user_id = $1
	RETURNING id
	`
	data := &ModifiedData{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(&data.ID)
	if err != nil {
		return nil, DetermineDBError(err, "resettoken_delete")
	}
	data.Timestamp = time.Now().UTC()
	return data, nil
}

func (m ResetTokenModel) DeleteByToken(ctx context.Context, token string) (*ModifiedData, error) {
	const query = `
	DELETE FROM reset_tokens 
	WHERE reset_token = $1
	RETURNING id
	`
	data := &ModifiedData{
		Timestamp: time.Now().UTC(),
	}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		token,
	).Scan(&data.ID)
	if err != nil {
		return nil, DetermineDBError(err, "resettoken_delete")
	}
	return data, nil
}
