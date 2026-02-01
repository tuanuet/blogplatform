-- User Profile Migration
-- Add profile fields to users table

-- Add profile columns to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS display_name VARCHAR(50),
ADD COLUMN IF NOT EXISTS bio TEXT,
ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500),
ADD COLUMN IF NOT EXISTS website VARCHAR(255),
ADD COLUMN IF NOT EXISTS location VARCHAR(100),
ADD COLUMN IF NOT EXISTS twitter_handle VARCHAR(50),
ADD COLUMN IF NOT EXISTS github_handle VARCHAR(50),
ADD COLUMN IF NOT EXISTS linkedin_url VARCHAR(255);

-- Create index for public profile lookups
CREATE INDEX IF NOT EXISTS idx_users_display_name ON users(display_name);

-- Comment on columns
COMMENT ON COLUMN users.display_name IS 'Public display name (3-50 chars)';
COMMENT ON COLUMN users.bio IS 'User biography (max 500 chars)';
COMMENT ON COLUMN users.avatar_url IS 'Avatar image URL (uploaded to server)';
COMMENT ON COLUMN users.website IS 'Personal website URL';
COMMENT ON COLUMN users.location IS 'User location (City, Country)';
COMMENT ON COLUMN users.twitter_handle IS 'Twitter username without @';
COMMENT ON COLUMN users.github_handle IS 'GitHub username';
COMMENT ON COLUMN users.linkedin_url IS 'LinkedIn profile URL';
