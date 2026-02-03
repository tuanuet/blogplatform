-- Rollback: Drop SePay-related tables and revert subscriptions changes

-- Drop triggers first
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;

-- Drop indexes
DROP INDEX IF EXISTS idx_subscriptions_expires_at;
DROP INDEX IF EXISTS idx_transactions_user_id;
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_transactions_sepay_id;
DROP INDEX IF EXISTS idx_transactions_reference_code;
DROP INDEX IF EXISTS idx_transactions_order_id;
DROP INDEX IF EXISTS idx_user_series_purchases_series_id;

-- Drop tables in reverse order of creation (respecting dependencies)
DROP TABLE IF EXISTS user_series_purchases;
DROP TABLE IF EXISTS transactions;

-- Revert subscriptions table changes
ALTER TABLE subscriptions
DROP COLUMN IF EXISTS expires_at,
DROP COLUMN IF EXISTS tier,
DROP COLUMN IF EXISTS updated_at;
