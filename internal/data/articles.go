package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Article struct {
	ID          string    `json:"id"`
	AuthorID    string    `json:"author_id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	Category    string    `json:"category"`
	TagIDs      []int     `json:"tag_ids,omitempty"`
	CategoryID  string    `json:"category_id"`
	Content     string    `json:"content"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"last_updated_at"`
	PublishedAt time.Time `json:"published_at,omitempty"`
}

type ArticleModel struct {
	DB *sql.DB
}

func (m *ArticleModel) Create(ctx context.Context, article *Article) (*Article, error) {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, DetermineDBError(err, "article_create")
	}
	defer tx.Rollback()

	const query = `
	INSERT INTO articles (author_id, title, content, status)
	VALUES ($1, $2, $3, $4)
	RETURNING id, author_id, title, content, status
	`
	newArticle := &Article{}
	err = tx.QueryRowContext(
		ctx,
		query,
		article.AuthorID,
		article.Title,
		article.Content,
		article.Status,
	).Scan(
		&newArticle.ID,
		&newArticle.AuthorID,
		&newArticle.Title,
		&newArticle.Content,
		&newArticle.Status,
	)
	if err != nil {
		return nil, DetermineDBError(err, "article_create")
	}
	if len(article.TagIDs) > 0 {
		err = m.attachTags(ctx, tx, newArticle.ID, article.TagIDs)
		if err != nil {
			return nil, DetermineDBError(err, "article_attachtags")
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, DetermineDBError(err, "article_create")
	}
	if len(article.Tags) > 0 {
		newArticle.Tags = article.Tags
	}
	return newArticle, nil
}

func (m *ArticleModel) attachTags(ctx context.Context, tx *sql.Tx, articleID string, tagIDs []int) error {
	const query = `
	INSERT INTO article_tags (article_id, tag_id)
	VALUES ($1, $2)
	`
	for _, tagID := range tagIDs {
		_, err := tx.ExecContext(ctx, query, articleID, tagID)
		if err != nil {
			return errors.New("error occured in tag attachment")
		}
	}
	return nil
}

func (m *ArticleModel) GetByID(id string) (*Article, error) {
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

//func (m ArticleModel) GetAll(filters Filters) ([]Article, Metadata, error) {
//	q := `
//		SELECT COUNT(*) OVER() as total_count, a.id, a.author_id, a.title, a.content, a.status, a.created_at, a.updated_at, a.published_at,
//		c.name AS category
//		FROM articles a
//		LEFT JOIN article_categories ac ON ac.article_id = a.id
//		LEFT JOIN categories c ON ac.category_id = c.id
//		`
//
//	if filters.Search != "" {
//		q += `
//		AND (
//				to_tsvector('simple', a.title) @@ plainto_tsquery('simple', $1)
//                OR to_tsvector('simple', a.content) @@ plainto_tsquery('simple', $1)
//		)
//		`
//	}
//
//	if filters.Status != "" {
//		q += `
//		AND a.status = $2
//		`
//	}
//
//	if filters.Category != "" {
//		q += `
//		AND c.name = $3
//		`
//	}
//
//	if len(filters.Tags) > 0 {
//		q += `
//		AND t.name IN ($4)
//		`
//	}
//
//	q += `
//        ORDER BY a.created_at DESC
//        LIMIT $5 OFFSET $6
//    `
//	args := []any{filters.Search, filters.Status, filters.Category, filters.limit(), filters.offset()}
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	rows, err := m.DB.QueryContext(ctx, q, args...)
//	if err != nil {
//		return nil, Metadata{}, err
//	}
//	defer rows.Close()
//	articles := make([]Article, 0)
//	var totalCount int
//	if err = rows.Scan() {
//		var article Article
//
//	}
//}

func (m *ArticleModel) Update(article *Article) (*Article, error) {
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

func (m *ArticleModel) Delete(id string) error {
	q := `DELETE FROM articles WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, q, id)
	if err != nil {
		return DetermineDBError(err, "article_delete")
	}
	return nil
}
