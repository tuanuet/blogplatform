-- Drop triggers
DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;

-- Drop tables in reverse order (respecting dependencies)
DROP TABLE IF EXISTS user_device_tokens;
DROP TABLE IF EXISTS notification_preferences;
DROP TABLE IF EXISTS notifications;
