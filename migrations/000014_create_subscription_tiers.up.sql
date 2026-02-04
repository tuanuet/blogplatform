-- Migration: Create subscription tier-related tables
-- Description: Creates subscription_plans and tag_tier_mappings tables, adds tier column to subscriptions table

-- Create subscription_plans table
CREATE TABLE IF NOT EXISTS subscription_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tier VARCHAR(20) NOT NULL CHECK (tier IN ('FREE', 'BRONZE', 'SILVER', 'GOLD')),
    price DECIMAL(15,2) NOT NULL CHECK (price >= 0),
    duration_days INT NOT NULL DEFAULT 30 CHECK (duration_days > 0),
    name VARCHAR(100),
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    UNIQUE(author_id, tier)
);

-- Create tag_tier_mappings table
CREATE TABLE IF NOT EXISTS tag_tier_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    required_tier VARCHAR(20) NOT NULL CHECK (required_tier IN ('FREE', 'BRONZE', 'SILVER', 'GOLD')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(author_id, tag_id)
);

-- Create indexes for subscription_plans
CREATE INDEX IF NOT EXISTS idx_subscription_plans_author_tier ON subscription_plans(author_id, tier);
CREATE INDEX IF NOT EXISTS idx_subscription_plans_author_id ON subscription_plans(author_id);
CREATE INDEX IF NOT EXISTS idx_subscription_plans_is_active ON subscription_plans(is_active);
CREATE INDEX IF NOT EXISTS idx_subscription_plans_deleted_at ON subscription_plans(deleted_at);

-- Create indexes for tag_tier_mappings
CREATE INDEX IF NOT EXISTS idx_tag_tier_mappings_author_tag ON tag_tier_mappings(author_id, tag_id);
CREATE INDEX IF NOT EXISTS idx_tag_tier_mappings_author_id ON tag_tier_mappings(author_id);
CREATE INDEX IF NOT EXISTS idx_tag_tier_mappings_tag_id ON tag_tier_mappings(tag_id);

-- Create trigger for subscription_plans updated_at
CREATE TRIGGER update_subscription_plans_updated_at
    BEFORE UPDATE ON subscription_plans
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create trigger for tag_tier_mappings updated_at
CREATE TRIGGER update_tag_tier_mappings_updated_at
    BEFORE UPDATE ON tag_tier_mappings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
