-- =============================================
-- User Reading History Migration
-- =============================================

CREATE TABLE IF NOT EXISTS user_reading_history (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blog_id UUID NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    last_read_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, blog_id)
);

-- Index for retrieving a user's history sorted by time
CREATE INDEX idx_user_reading_history_user_last_read ON user_reading_history(user_id, last_read_at DESC);
