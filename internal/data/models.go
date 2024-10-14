package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users UserModel
}

func AppModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}
