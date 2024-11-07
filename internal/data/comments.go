package data

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ArticleID string    `json:"article_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Create(ctx context.Context, comment *Comment) (*Comment, error) {
	const query = `
	INSERT INTO comments (user_id, article_id, content)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, article_id, content, created_at
	`
	newComment := &Comment{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		comment.UserID,
		comment.ArticleID,
		comment.Content,
	).Scan(
		&newComment.ID,
		&newComment.UserID,
		&newComment.ArticleID,
		&newComment.Content,
		&newComment.CreatedAt,
	)
	if err != nil {
		return nil, DetermineDBError(err, "comment_create")
	}
	return newComment, nil
}

func (m *CommentModel) Delete(ctx context.Context, id string) (*ModifiedData, error) {
	const query = `
	DELETE FROM comments 
	WHERE id = $1
	RETURNING id
	`
	data := &ModifiedData{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&data.ID,
	)
	if err != nil {
		return nil, DetermineDBError(err, "comment_delete")
	}
	data.Timestamp = time.Now().UTC()
	return data, nil
}
