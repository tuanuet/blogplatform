-- Rollback: Drop subscription tier-related tables

-- Drop triggers first
DROP TRIGGER IF EXISTS update_tag_tier_mappings_updated_at ON tag_tier_mappings;
DROP TRIGGER IF EXISTS update_subscription_plans_updated_at ON subscription_plans;

-- Drop indexes for tag_tier_mappings
DROP INDEX IF EXISTS idx_tag_tier_mappings_tag_id;
DROP INDEX IF EXISTS idx_tag_tier_mappings_author_id;
DROP INDEX IF EXISTS idx_tag_tier_mappings_author_tag;

-- Drop indexes for subscription_plans
DROP INDEX IF EXISTS idx_subscription_plans_is_active;
DROP INDEX IF EXISTS idx_subscription_plans_author_id;
DROP INDEX IF EXISTS idx_subscription_plans_author_tier;

-- Drop tables in reverse order of creation (respecting dependencies)
DROP TABLE IF EXISTS tag_tier_mappings;
DROP TABLE IF EXISTS subscription_plans;
