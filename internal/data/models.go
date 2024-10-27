package data

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrEditConflict        = errors.New("edit conflict")
	ErrDuplicateKey        = errors.New("duplicate key value violates unique constraint")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrCheckConstraint     = errors.New("check constraint violation")
	ErrInvalidInput        = errors.New("invalid input syntax")
	ErrConnectionFailed    = errors.New("connection failed")
)

type Models struct {
	Users       UserModel
	ResetTokens ResetTokenModel
}

type DBError struct {
	Err       error
	Operation string
	Detail    string
}

func (e DBError) Error() string {
	return fmt.Sprintf("database error during %s: %v - %s", e.Operation, e.Err, e.Detail)
}

func (e DBError) Unwrap() error {
	return e.Err
}

func DetermineDBError(err error, operation string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return &DBError{
			Err:       ErrRecordNotFound,
			Operation: operation,
			Detail:    "requested record does not exist",
		}
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return categorizePostgresError(pqErr, operation)
	}

	if errors.Is(err, sql.ErrConnDone) || errors.Is(err, sql.ErrTxDone) {
		return &DBError{
			Err:       ErrConnectionFailed,
			Operation: operation,
			Detail:    "database connection is closed",
		}
	}

	return &DBError{
		Err:       err,
		Operation: operation,
		Detail:    "uncategorized database error",
	}
}

func categorizePostgresError(pqErr *pq.Error, operation string) error {
	var err error
	detail := pqErr.Detail
	switch pqErr.Code {
	case "23505": // unique_violation
		err = ErrDuplicateKey
	case "23503": // foreign_key_violation
		err = ErrForeignKeyViolation
	case "23514": // check_violation
		err = ErrCheckConstraint
	case "22P02": // invalid_text_representation
		err = ErrInvalidInput
	case "23000": // integrity_constraint_violation
		err = ErrEditConflict
	default:
		// Log unknown error codes for future categorization
		err = fmt.Errorf("unhandled postgres error: %s", pqErr.Code)
	}

	return &DBError{
		Err:       err,
		Operation: operation,
		Detail:    detail,
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:       UserModel{DB: db},
		ResetTokens: ResetTokenModel{DB: db},
	}
}
