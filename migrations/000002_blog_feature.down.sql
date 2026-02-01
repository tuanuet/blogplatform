-- Drop triggers
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS update_blogs_updated_at ON blogs;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP TRIGGER IF EXISTS update_tags_updated_at ON tags;

-- Drop tables in correct order (respecting foreign keys)
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS blog_tags;
DROP TABLE IF EXISTS blogs;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;

-- Drop enum types
DROP TYPE IF EXISTS blog_visibility;
DROP TYPE IF EXISTS blog_status;
