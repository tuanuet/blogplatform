-- Create user velocity scores table
CREATE TABLE IF NOT EXISTS user_velocity_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    follower_count INTEGER NOT NULL DEFAULT 0,
    follower_growth_rate DECIMAL(10, 4) NOT NULL DEFAULT 0,
    blog_post_velocity DECIMAL(10, 4) NOT NULL DEFAULT 0,
    composite_score DECIMAL(10, 4) NOT NULL DEFAULT 0,
    rank_position INTEGER,
    calculation_date TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for user velocity scores
CREATE INDEX idx_velocity_scores_composite ON user_velocity_scores(composite_score DESC);
CREATE INDEX idx_velocity_scores_rank ON user_velocity_scores(rank_position) WHERE rank_position IS NOT NULL;
CREATE INDEX idx_velocity_scores_user_id ON user_velocity_scores(user_id);

-- Create user ranking history table
CREATE TABLE IF NOT EXISTS user_ranking_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rank_position INTEGER NOT NULL,
    composite_score DECIMAL(10, 4) NOT NULL,
    follower_count INTEGER NOT NULL DEFAULT 0,
    recorded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for ranking history
CREATE INDEX idx_ranking_history_user_date ON user_ranking_history(user_id, recorded_at DESC);
CREATE INDEX idx_ranking_history_recorded_at ON user_ranking_history(recorded_at);

-- Create follower snapshots table for velocity calculation
CREATE TABLE IF NOT EXISTS user_follower_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    follower_count INTEGER NOT NULL,
    snapshot_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    UNIQUE(user_id, snapshot_date)
);

-- Indexes for follower snapshots
CREATE INDEX idx_follower_snapshots_user_date ON user_follower_snapshots(user_id, snapshot_date);
CREATE INDEX idx_follower_snapshots_date ON user_follower_snapshots(snapshot_date);

-- Apply updated_at trigger to new tables
CREATE TRIGGER update_user_velocity_scores_updated_at
    BEFORE UPDATE ON user_velocity_scores
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
