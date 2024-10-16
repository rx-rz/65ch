package data

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func NilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("nil"), nil
	}
	return json.Marshal(ns.String)
}

func (ns NullString) UnmarshalJSON(b []byte) error {
	var s *string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}
