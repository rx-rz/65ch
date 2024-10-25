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
	IsUsed     bool      `json:"is_used"`
	Expiration time.Time `json:"expiration"`
}

func (m ResetTokenModel) Create(resetToken *ResetToken) error {
	q := `
	INSERT INTO  reset_tokens (user_id, reset_token, expiration) VALUES ($1, $2, $3)
	`
	args := []any{resetToken.UserID, resetToken.ResetToken, resetToken.Expiration}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, q, args...)
	return determineDBError(err)

}

func (m ResetTokenModel) Update(resetToken *ResetToken) error {
	return nil
	//	q := `
	//	UPDATE reset_tokens SET
	//`
}
