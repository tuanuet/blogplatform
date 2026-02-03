-- Migration: Create SePay-related tables
-- Description: Creates transactions and user_series_purchases tables for payment handling, and adds expires_at, tier, updated_at columns to subscriptions table

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(20,2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'VND',
    provider VARCHAR(20) NOT NULL,
    gateway VARCHAR(50),
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    target_id UUID,
    plan_id VARCHAR(50),
    content TEXT,
    sepay_id VARCHAR(100) UNIQUE,
    reference_code VARCHAR(100),
    order_id VARCHAR(100),
    webhook_payload JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for transactions table
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE UNIQUE INDEX IF NOT EXISTS idx_transactions_sepay_id ON transactions(sepay_id);
CREATE INDEX IF NOT EXISTS idx_transactions_reference_code ON transactions(reference_code);
CREATE INDEX IF NOT EXISTS idx_transactions_order_id ON transactions(order_id);

-- Create trigger for transactions updated_at
CREATE TRIGGER update_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create user_series_purchases table
CREATE TABLE IF NOT EXISTS user_series_purchases (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    series_id UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    amount DECIMAL(20,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, series_id)
);

-- Create index for user_series_purchases
CREATE INDEX IF NOT EXISTS idx_user_series_purchases_series_id ON user_series_purchases(series_id);

-- Alter subscriptions table to add new columns
ALTER TABLE subscriptions
ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS tier VARCHAR(20) DEFAULT 'FREE',
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT NOW();

-- Create index for subscriptions.expires_at
CREATE INDEX IF NOT EXISTS idx_subscriptions_expires_at ON subscriptions(expires_at);

-- Create trigger for subscriptions updated_at
CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
