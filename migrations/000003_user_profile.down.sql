-- Rollback User Profile Migration
-- Remove profile fields from users table

ALTER TABLE users
DROP COLUMN IF EXISTS display_name,
DROP COLUMN IF EXISTS bio,
DROP COLUMN IF EXISTS avatar_url,
DROP COLUMN IF EXISTS website,
DROP COLUMN IF EXISTS location,
DROP COLUMN IF EXISTS twitter_handle,
DROP COLUMN IF EXISTS github_handle,
DROP COLUMN IF EXISTS linkedin_url;

DROP INDEX IF EXISTS idx_users_display_name;
