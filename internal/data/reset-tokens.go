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

func (m ResetTokenModel) Create(resetToken *ResetToken) error {
	q := `
	INSERT INTO  reset_tokens (user_id, reset_token, expiration) VALUES ($1, $2, $3)
	`
	args := []any{resetToken.UserID, resetToken.ResetToken, resetToken.Expiration}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, q, args...)
	return DetermineDBError(err, "resettoken_create")

}

func (m ResetTokenModel) GetByUserID(userId string) (*ResetToken, error) {
	var resetToken ResetToken
	q := `SELECT id, user_id, reset_token, expiration FROM reset_tokens WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, userId).Scan(&resetToken.ID, &resetToken.UserID, &resetToken.ResetToken, &resetToken.Expiration)
	if err != nil {
		return nil, DetermineDBError(err, "resettoken_getbyuserid")
	}
	return &resetToken, nil
}

func (m ResetTokenModel) GetByToken(token string) (*ResetToken, error) {
	var resetToken ResetToken
	q := `SELECT id, user_id, reset_token, expiration FROM reset_tokens WHERE reset_token = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, token).Scan(&resetToken.ID, &resetToken.UserID, &resetToken.ResetToken, &resetToken.Expiration)
	if err != nil {
		return nil, DetermineDBError(err, "resettoken_getbytoken")
	}
	return &resetToken, nil
}

func (m ResetTokenModel) Update(resetToken *ResetToken) error {
	q := `UPDATE reset_tokens SET reset_token = $1, expiration = $2 WHERE user_id = $3`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	args := []any{resetToken.ResetToken, resetToken.Expiration, resetToken.UserID}
	_, err := m.DB.ExecContext(ctx, q, args...)
	if err != nil {
		return DetermineDBError(err, "resettoken_update")
	}
	return nil
}

func (m ResetTokenModel) Delete(userId string) error {
	q := `
	DELETE FROM reset_tokens WHERE user_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, q, userId)
	return DetermineDBError(err, "resettoken_delete")
}
