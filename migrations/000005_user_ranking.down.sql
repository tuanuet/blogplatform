-- Drop triggers
DROP TRIGGER IF EXISTS update_user_velocity_scores_updated_at ON user_velocity_scores;

-- Drop tables
DROP TABLE IF EXISTS user_follower_snapshots;
DROP TABLE IF EXISTS user_ranking_history;
DROP TABLE IF EXISTS user_velocity_scores;
