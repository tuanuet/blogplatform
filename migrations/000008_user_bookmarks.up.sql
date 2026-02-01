-- =============================================
-- User Bookmarks Feature Migration
-- =============================================

-- User Bookmarks Join Table
CREATE TABLE IF NOT EXISTS user_bookmarks (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blog_id UUID NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, blog_id)
);

-- Indexes for performance
-- Primary key covers (user_id, blog_id), so we need one for blog_id for reverse lookups (e.g. count bookmarks per blog)
CREATE INDEX idx_user_bookmarks_blog_id ON user_bookmarks(blog_id);
-- Index for created_at might be useful for "recently bookmarked" queries
CREATE INDEX idx_user_bookmarks_created_at ON user_bookmarks(created_at);
