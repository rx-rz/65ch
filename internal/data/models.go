package data

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users       UserModel
	ResetTokens ResetTokenModel
}

func determineDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRecordNotFound
	}
	var dbError *pq.Error
	errors.As(err, &dbError)
	switch dbError.Code {
	case "23505":
		return ErrEditConflict
	}
	return nil
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:       UserModel{DB: db},
		ResetTokens: ResetTokenModel{DB: db},
	}
}
