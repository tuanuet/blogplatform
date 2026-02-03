ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified_at TIMESTAMP;

CREATE TABLE IF NOT EXISTS social_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_social_provider_id ON social_accounts(provider, provider_id);
CREATE INDEX IF NOT EXISTS idx_social_user_id ON social_accounts(user_id);

CREATE TRIGGER update_social_accounts_updated_at
    BEFORE UPDATE ON social_accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
