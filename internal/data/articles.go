package data

import (
	"context"
	"database/sql"
	"time"
)

type Article struct {
	ID          string    `json:"id"`
	AuthorID    string    `json:"author_id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"last_updated_at"`
	PublishedAt time.Time `json:"published_at"`
}

type ArticleModel struct {
	DB *sql.DB
}

func (m ArticleModel) Create(article *Article) error {
	q := `INSERT INTO articles (author_id, title, content) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	args := []any{article.AuthorID, article.Title, article.Content}
	_, err := m.DB.ExecContext(ctx, q, args)
	if err != nil {
		return DetermineDBError(err, "article_create")
	}
	return nil
}

func (m ArticleModel) GetByID(id string) (*Article, error) {
	var article *Article
	q := `SELECT id, author_id, title, content, created_at, updated_at, published_at FROM articles WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, id).Scan(&article.ID, &article.AuthorID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt, &article.PublishedAt)
	if err != nil {
		return nil, DetermineDBError(err, "article_getbyid")
	}
	return article, nil
}

func (m ArticleModel) GetAllPublishedArticles() ([]Article, error) {
	q := `SELECT count(*) OVER(), id, author_id, title, content, created_at, updated_at, published_at FROM articles WHERE status = 'published'`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q)
}

func (m ArticleModel) Update(article *Article) (*Article, error) {
	var articleDetails *Article
	q := `UPDATE articles SET title = $1, content = $2, updated_at = $3, published_at = $4 WHERE id = $5 RETURNING id, updated_at`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	args := []any{article.AuthorID, article.Title, article.UpdatedAt, article.PublishedAt, article.ID}
	err := m.DB.QueryRowContext(ctx, q, args...).Scan(&articleDetails.ID, &articleDetails.UpdatedAt)
	if err != nil {
		return nil, DetermineDBError(err, "article_update")
	}
	return articleDetails, nil
}

func (m ArticleModel) Delete(id string) error {
	q := `DELETE FROM articles WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, q, id)
	if err != nil {
		return DetermineDBError(err, "article_delete")
	}
	return nil
}
