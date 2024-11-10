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
	CategoryID  int       `json:"category_id"`
	Content     string    `json:"content"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"last_updated_at"`
	PublishedAt time.Time `json:"published_at,omitempty"`
}

type SavedArticle struct {
	UserID    string    `json:"user_id"`
	ArticleID string    `json:"article_id"`
	SavedAt   time.Time `json:"saved_at"`
}

type LikedArticle struct {
	UserID    string `json:"user_id"`
	ArticleID string `json:"article_id"`
	LikedAt   string `json:"liked_at"`
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
	INSERT INTO articles (author_id, title, content, status, published_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, author_id, title, content, status
	`
	newArticle := &Article{}
	publishTimestamp := time.Now().UTC()

	err = tx.QueryRowContext(
		ctx,
		query,
		article.AuthorID,
		article.Title,
		article.Content,
		article.Status,
		publishTimestamp,
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
	const deleteQuery = `
	DELETE FROM article_tags
	WHERE article_id = $1
	`

	_, err := tx.ExecContext(ctx, deleteQuery, articleID)
	if err != nil {
		return errors.New("error occured in tag attachment")
	}

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

func (m *ArticleModel) GetByID(ctx context.Context, id string) (*Article, error) {
	q := `
	SELECT id, author_id, title, content, created_at, updated_at, published_at 
	FROM articles 
	WHERE id = $1`

	article := &Article{}

	err := m.DB.QueryRowContext(
		ctx,
		q,
		id,
	).Scan(
		&article.ID,
		&article.AuthorID,
		&article.Title,
		&article.Content,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.PublishedAt,
	)
	if err != nil {
		return nil, DetermineDBError(err, "article_getbyid")
	}
	return article, nil
}

func (m *ArticleModel) Update(ctx context.Context, article *Article) (*ModifiedData, error) {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, DetermineDBError(err, "article_create")
	}
	defer tx.Rollback()

	const query = `
	UPDATE articles 
	SET title = $1, 
		content = $2, 
		status = $3,
		category_id = $4,
		updated_at = $5,
		published_at = $6 
	WHERE id = $7
	RETURNING id
	`

	updateTimestamp := time.Now().UTC()
	data := &ModifiedData{}

	err = tx.QueryRowContext(
		ctx,
		query,
		article.Title,
		article.Content,
		article.Status,
		article.CategoryID,
		updateTimestamp,
		article.PublishedAt,
		article.ID,
	).Scan(
		&data.ID,
	)

	if err != nil {
		return nil, DetermineDBError(err, "article_update")
	}
	if len(article.TagIDs) > 0 {
		err = m.attachTags(ctx, tx, data.ID, article.TagIDs)
		if err != nil {
			return nil, DetermineDBError(err, "article_attachtags")
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, DetermineDBError(err, "article_create")
	}
	data.Timestamp = updateTimestamp

	return data, nil
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

func (m *ArticleModel) Save(ctx context.Context, article *SavedArticle) (*SavedArticle, error) {
	const query = `
	INSERT INTO saved_articles 
	VALUES ($1, $2)
	RETURNING  user_id, article_id
	`
	savedArticle := &SavedArticle{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		article,
		article.UserID,
		article.ArticleID,
	).Scan(
		&savedArticle.UserID,
		&savedArticle.ArticleID,
	)
	if err != nil {
		return nil, DetermineDBError(err, "article_save")
	}
	return savedArticle, nil
}

func (m *ArticleModel) Unsave(ctx context.Context, userID, articleID string) (*ModifiedData, error) {
	const query = `
	DELETE FROM saved_articles
	WHERE user_id = $1 AND article_id = $2
	RETURNING (user_id, article_id)
	`
	data := &ModifiedData{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		userID,
		articleID,
	).Scan(
		&data.ID,
	)
	if err != nil {
		return nil, DetermineDBError(err, "article_unsave")
	}
	data.Timestamp = time.Now().UTC()
	return data, nil
}

func (m *ArticleModel) Like(ctx context.Context, userID, articleID string) (*LikedArticle, error) {
	const query = `
	INSERT INTO liked_articles
	VALUES ($1, $2)
	RETURNING user_id, article_id, liked_at
	`
	likedArticle := &LikedArticle{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		userID,
		articleID,
	).Scan(
		&likedArticle.UserID,
		&likedArticle.ArticleID,
		&likedArticle.LikedAt,
	)
	if err != nil {
		return nil, DetermineDBError(err, "article_like")
	}
	return likedArticle, nil
}

func (m *ArticleModel) Unlike(ctx context.Context, userID, articleID string) (*ModifiedData, error) {
	const query = `
	DELETE FROM liked_articles
	WHERE user_id = $1 
	AND article_id = $2
	RETURNING id
	`
	data := &ModifiedData{}
	err := m.DB.QueryRowContext(
		ctx,
		query,
		userID,
		articleID,
	).Scan(
		&data.ID,
	)
	data.Timestamp = time.Now().UTC()
	if err != nil {
		return nil, DetermineDBError(err, "article_unlike")

	}
	return data, nil
}
