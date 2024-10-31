package data

import (
	"context"
	"database/sql"
	"strconv"
	"time"
)

type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"string"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"last_updated_at"`
}

type TagModel struct {
	DB *sql.DB
}

func (m TagModel) Create(name string) (*Tag, error) {
	var tag Tag
	q := `
	INSERT INTO tags (name) VALUES ($1) RETURNING id, name, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, name).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, DetermineDBError(err, "tag_create")
	}
	return &tag, nil
}

func (m TagModel) GetAll() ([]*Tag, error) {
	q := `
	SELECT id, name FROM tags	ORDER BY created_at DESC 
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, DetermineDBError(err, "tag_getall")
	}
	defer rows.Close()
	var tags []*Tag
	for rows.Next() {
		var c Tag
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, DetermineDBError(err, "tag_getall")
		}
		tags = append(tags, &c)
	}
	return tags, nil
}

func (m TagModel) GetByName(name string) (*Tag, error) {
	var tag Tag
	q := `
	SELECT id, name, created_at, updated_at FROM tags WHERE name = $1 ORDER BY id DESC
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, name).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, DetermineDBError(err, "tag_getbyname")
	}
	return &tag, nil
}

func (m TagModel) GetByID(id string) (*Tag, error) {
	var tag Tag
	q := `
	SELECT id, name, created_at, updated_at FROM tags WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, id).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, DetermineDBError(err, "tag_getbyid")
	}
	return &tag, nil
}

func (m TagModel) UpdateName(tag Tag) (ModifiedData, error) {
	q := `
	UPDATE tags SET name = $1, updated_at = current_timestamp WHERE id = $2 
	`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	args := []any{tag.Name, tag.ID}
	_, err := m.DB.ExecContext(ctx, q, args...)
	if err != nil {
		return ModifiedData{}, DetermineDBError(err, "tag_updatename")
	}
	return ModifiedData{
		ID:        strconv.Itoa(tag.ID),
		Timestamp: time.Now(),
	}, nil
}

func (m TagModel) Delete(id int) (ModifiedData, error) {
	q := `DELETE FROM tags WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, q, id)
	if err != nil {
		return ModifiedData{}, DetermineDBError(err, "tag_deletebyname")
	}
	return ModifiedData{
		ID:        strconv.Itoa(id),
		Timestamp: time.Now(),
	}, nil

}
