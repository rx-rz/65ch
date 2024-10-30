package data

import "database/sql"

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"string"`
}

type TagModel struct {
	DB *sql.DB
}
