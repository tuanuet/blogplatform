-- Rollback: Drop blog_versions and blog_version_tags tables

-- Drop tables in reverse order of creation (respecting dependencies)
DROP TABLE IF EXISTS blog_version_tags;
DROP TABLE IF EXISTS blog_versions;

-- Drop indexes (PostgreSQL drops indexes automatically when table is dropped,
-- but we include them for explicitness and other database compatibility)
DROP INDEX IF EXISTS idx_blog_versions_blog_id_created_at;
DROP INDEX IF EXISTS idx_blog_versions_blog_id;
