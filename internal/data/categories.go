package data

import (
	"context"
	"database/sql"
	"time"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryModel struct {
	DB *sql.DB
}

func (m CategoryModel) Create(name string) (*Category, error) {
	var category Category
	q := `
	INSERT INTO categories (name) VALUES ($1) RETURNING id, name
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, name).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, DetermineDBError(err, "category_create")
	}
}

func (m CategoryModel) GetAll() {

}

func (m CategoryModel) GetByName(identifier map[string]string) {

}

func (m CategoryModel) GetByID(id string) {

}

func (m CategoryModel) UpdateName(name string) {

}

func (m CategoryModel) DeleteByName(name string) {

}

func (m CategoryModel) DeleteByID(id string) {

}
