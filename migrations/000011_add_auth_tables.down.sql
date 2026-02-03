DROP TRIGGER IF EXISTS update_social_accounts_updated_at ON social_accounts;
DROP TABLE IF EXISTS social_accounts;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified_at;
