CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE users ALTER COLUMN profile_picture_url SET DATA TYPE TEXT;
-- Users Table
CREATE TABLE users
(
    id                  UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    email               VARCHAR(255) UNIQUE NOT NULL,
    password_hash       bytea        NOT NULL,
    first_name          VARCHAR(255)        NOT NULL,
    last_name           VARCHAR(255)        NOT NULL,
    activated           BOOLEAN                  DEFAULT FALSE,
    bio                 TEXT                     DEFAULT 'bio',
    profile_picture_url TEXT             DEFAULT 'https://placehold.co/400?text=U',
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for email to speed up lookups
CREATE INDEX idx_users_email ON users (email);

-- Articles Table
CREATE TABLE articles
(
    id           UUID PRIMARY KEY                                                 DEFAULT uuid_generate_v4(),
    author_id    UUID REFERENCES users (id),
    title        VARCHAR(255) NOT NULL,
    content      TEXT,
    status       VARCHAR(20) CHECK (status IN ('draft', 'published', 'archived')) DEFAULT 'draft',
    created_at   TIMESTAMP WITH TIME ZONE                                         DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE                                         DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP WITH TIME ZONE
);

-- Index for author_id and status
CREATE INDEX idx_articles_author_id ON articles (author_id);
CREATE INDEX idx_articles_status ON articles (status);

-- Tags Table
CREATE TABLE tags
(
    id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Index for tag name
CREATE INDEX idx_tags_name ON tags (name);

-- Article Tags (Many-to-Many Relationship)
CREATE TABLE article_tags
(
    article_id UUID REFERENCES articles (id),
    tag_id     UUID REFERENCES tags (id),
    PRIMARY KEY (article_id, tag_id)
);

-- Comments Table
CREATE TABLE comments
(
    id         UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    article_id UUID REFERENCES articles (id),
    user_id    UUID REFERENCES users (id),
    content    TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for article_id and user_id in comments
CREATE INDEX idx_comments_article_id ON comments (article_id);
CREATE INDEX idx_comments_user_id ON comments (user_id);

-- Likes Table
CREATE TABLE likes
(
    id         UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    article_id UUID REFERENCES articles (id),
    user_id    UUID REFERENCES users (id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (article_id, user_id)
);

-- Index for article_id and user_id in likes
CREATE INDEX idx_likes_article_id ON likes (article_id);
CREATE INDEX idx_likes_user_id ON likes (user_id);

-- Followers Table
CREATE TABLE followers
(
    follower_id UUID REFERENCES users (id),
    followed_id UUID REFERENCES users (id),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, followed_id)
);

-- Saved Articles Table
CREATE TABLE saved_articles
(
    user_id    UUID REFERENCES users (id),
    article_id UUID REFERENCES articles (id),
    saved_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, article_id)
);

-- Index for user_id and article_id in saved_articles
CREATE INDEX idx_saved_articles_user_id ON saved_articles (user_id);
CREATE INDEX idx_saved_articles_article_id ON saved_articles (article_id);

-- Notifications Table
CREATE TABLE notifications
(
    id         UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    user_id    UUID REFERENCES users (id),
    type       VARCHAR(50) NOT NULL,
    content    TEXT        NOT NULL,
    is_read    BOOLEAN                  DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for user_id in notifications
CREATE INDEX idx_notifications_user_id ON notifications (user_id);

-- Categories Table
CREATE TABLE categories
(
    id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Index for category name
CREATE INDEX idx_categories_name ON categories (name);

-- Article Categories (Many-to-Many Relationship)
CREATE TABLE article_categories
(
    article_id  UUID REFERENCES articles (id),
    category_id UUID REFERENCES categories (id),
    PRIMARY KEY (article_id, category_id)
);

-- Article Statistics Table
CREATE TABLE article_statistics
(
    article_id     UUID REFERENCES articles (id),
    views_count    INTEGER                  DEFAULT 0,
    likes_count    INTEGER                  DEFAULT 0,
    comments_count INTEGER                  DEFAULT 0,
    last_updated   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (article_id)
);

