package data

import (
	"context"
	"database/sql"
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

func (m *TagModel) Create(ctx context.Context, name string) (*Tag, error) {
	const query = `
	INSERT INTO tags (name)
	VALUES ($1) 
	RETURNING id, name, created_at, updated_at
	`
	newTag := &Tag{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		name,
	).Scan(
		&newTag.ID,
		&newTag.Name,
		&newTag.CreatedAt,
		&newTag.UpdatedAt,
	)
	if err != nil {
		return nil, DetermineDBError(err, "tag_create")
	}
	return newTag, nil
}

func (m *TagModel) GetAll(ctx context.Context) ([]*Tag, error) {
	const query = `
	SELECT id, name 
	FROM tags	
	ORDER BY created_at DESC 
	`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, DetermineDBError(err, "tag_getall")
	}
	defer rows.Close()

	var tags []*Tag

	for rows.Next() {
		var t Tag
		err := rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, DetermineDBError(err, "tag_getall")
		}
		tags = append(tags, &t)
	}
	return tags, nil
}

func (m *TagModel) GetByName(ctx context.Context, name string) (*Tag, error) {
	const query = `
	SELECT id, name, created_at, updated_at 
	FROM tags 
	WHERE name = $1 
	ORDER BY id DESC
	`
	tag := &Tag{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		name,
	).Scan(
		&tag.ID,
		&tag.Name,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)
	if err != nil {
		return nil, DetermineDBError(err, "tag_getbyname")
	}
	return tag, nil
}

func (m *TagModel) GetByID(ctx context.Context, id string) (*Tag, error) {
	const query = `
	SELECT id, name, created_at, updated_at 
	FROM tags 
	WHERE id = $1 
	ORDER BY id DESC
	`
	tag := &Tag{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&tag.ID,
		&tag.Name,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)
	if err != nil {
		return nil, DetermineDBError(err, "tag_getbyid")
	}
	return tag, nil
}

func (m *TagModel) UpdateName(ctx context.Context, tag *Tag) (*ModifiedData, error) {
	const query = `
	UPDATE tags SET name = $1, 
	updated_at = $2 
	WHERE id = $3
	`
	data := &ModifiedData{}
	updateTimestamp := time.Now().UTC()
	err := m.DB.QueryRowContext(
		ctx,
		query,
		tag.Name,
		updateTimestamp,
		tag.ID,
	).Scan(&data.ID)

	if err != nil {
		return nil, DetermineDBError(err, "tag_updatename")
	}
	data.Timestamp = updateTimestamp
	return data, nil
}

func (m *TagModel) Delete(ctx context.Context, id int) (*ModifiedData, error) {
	const query = `
	   DELETE FROM tags 
       WHERE id = $1
	`
	data := &ModifiedData{
		Timestamp: time.Now().UTC(),
	}
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&data.ID)
	if err != nil {
		return data, DetermineDBError(err, "tag_deletebyname")
	}
	return data, nil

}
