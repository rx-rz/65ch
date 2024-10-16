-- Remove Article Statistics Table
DROP TABLE IF EXISTS article_statistics;

-- Remove Article Categories (Many-to-Many Relationship) Table
DROP TABLE IF EXISTS article_categories;

-- Remove Categories Table
DROP TABLE IF EXISTS categories;

-- Remove Notifications Table
DROP TABLE IF EXISTS notifications;

-- Remove Saved Articles Table
DROP TABLE IF EXISTS saved_articles;

-- Remove Followers Table
DROP TABLE IF EXISTS followers;

-- Remove Likes Table
DROP TABLE IF EXISTS likes;

-- Remove Comments Table
DROP TABLE IF EXISTS comments;

-- Remove Article Tags (Many-to-Many Relationship) Table
DROP TABLE IF EXISTS article_tags;

-- Remove Tags Table
DROP TABLE IF EXISTS tags;

-- Remove Articles Table
DROP TABLE IF EXISTS articles;

-- Remove Users Table
DROP TABLE IF EXISTS users;

-- Drop all indexes created
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_articles_author_id;
DROP INDEX IF EXISTS idx_articles_status;
DROP INDEX IF EXISTS idx_tags_name;
DROP INDEX IF EXISTS idx_comments_article_id;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_likes_article_id;
DROP INDEX IF EXISTS idx_likes_user_id;
DROP INDEX IF EXISTS idx_saved_articles_user_id;
DROP INDEX IF EXISTS idx_saved_articles_article_id;
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP INDEX IF EXISTS idx_categories_name;

-- Remove the uuid-ossp extension
DROP EXTENSION IF EXISTS "uuid-ossp";
