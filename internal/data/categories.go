package data

import (
	"context"
	"database/sql"
	"strconv"
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
	return &category, nil
}

func (m CategoryModel) GetAll() ([]*Category, error) {
	q := `
	SELECT id, name FROM categories	
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, DetermineDBError(err, "category_getall")
	}
	defer rows.Close()
	var categories []*Category
	for rows.Next() {
		var c *Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, DetermineDBError(err, "category_getall")
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (m CategoryModel) GetByName(name string) (*Category, error) {
	var category *Category
	q := `
	SELECT id, name FROM categories WHERE name = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, name).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, DetermineDBError(err, "category_getbyname")
	}
	return category, nil
}

func (m CategoryModel) GetByID(id string) (*Category, error) {
	var category *Category
	q := `
	SELECT id, name FROM categories WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, id).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, DetermineDBError(err, "category_getbyid")
	}
	return category, nil
}

func (m CategoryModel) UpdateName(category Category) (ModifiedData, error) {
	q := `
	UPDATE categories SET name = $1 WHERE id = $2 
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	args := []any{category.Name, category.ID}
	_, err := m.DB.ExecContext(ctx, q, args...)
	if err != nil {
		return ModifiedData{}, DetermineDBError(err, "category_updatename")
	}
	return ModifiedData{
		ID:        strconv.Itoa(category.ID),
		Timestamp: time.Now(),
	}, nil
}

func (m CategoryModel) DeleteByName(name string) (ModifiedData, error) {
	var id string
	q := `DELETE FROM categories WHERE name = $1 RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, name).Scan(&id)
	if err != nil {
		return ModifiedData{}, DetermineDBError(err, "category_deletebyname")
	}
	return ModifiedData{
		ID:        id,
		Timestamp: time.Now(),
	}, nil

}

func (m CategoryModel) DeleteByID(id string) (ModifiedData, error) {
	q := `DELETE FROM categories WHERE name = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, q, id)
	if err != nil {
		return ModifiedData{}, DetermineDBError(err, "category_deletebyname")
	}
	return ModifiedData{
		ID:        id,
		Timestamp: time.Now(),
	}, nil

}
