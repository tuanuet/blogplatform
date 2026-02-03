package migrations

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	MIGRATION_NUMBER  = "000013"
	MIGRATION_NAME    = "create_sepay_tables"
	UP_FILE_PATTERN   = MIGRATION_NUMBER + "_" + MIGRATION_NAME + ".up.sql"
	DOWN_FILE_PATTERN = MIGRATION_NUMBER + "_" + MIGRATION_NAME + ".down.sql"
	MIGRATIONS_DIR    = "../../../migrations"
)

// readUpMigration reads and returns the up migration file content
func readUpMigration(t *testing.T) string {
	t.Helper()
	upPath := filepath.Join(MIGRATIONS_DIR, UP_FILE_PATTERN)
	content, err := os.ReadFile(upPath)
	require.NoError(t, err, "Should be able to read up migration file")
	return string(content)
}

// readDownMigration reads and returns the down migration file content
func readDownMigration(t *testing.T) string {
	t.Helper()
	downPath := filepath.Join(MIGRATIONS_DIR, DOWN_FILE_PATTERN)
	content, err := os.ReadFile(downPath)
	require.NoError(t, err, "Should be able to read down migration file")
	return string(content)
}

// TestMigrationFilesExist verifies that both up and down migration files exist
func TestMigrationFilesExist(t *testing.T) {
	// Check up file exists
	upPath := filepath.Join(MIGRATIONS_DIR, UP_FILE_PATTERN)
	_, err := os.Stat(upPath)
	if os.IsNotExist(err) {
		t.Fatalf("Migration up file does not exist: %s", upPath)
	}
	require.NoError(t, err, "Up migration file should exist")

	// Check down file exists
	downPath := filepath.Join(MIGRATIONS_DIR, DOWN_FILE_PATTERN)
	_, err = os.Stat(downPath)
	if os.IsNotExist(err) {
		t.Fatalf("Migration down file does not exist: %s", downPath)
	}
	require.NoError(t, err, "Down migration file should exist")
}

// TestUpMigrationContainsRequiredElements verifies the up migration contains all required SQL elements
func TestUpMigrationContainsRequiredElements(t *testing.T) {
	upSQL := readUpMigration(t)

	// Test: Contains CREATE TABLE transactions statement
	assert.Regexp(t, "(?i)CREATE TABLE( IF NOT EXISTS)? transactions", upSQL,
		"Up migration should contain CREATE TABLE transactions")

	// Test: transactions table has sepay_id column
	assert.Contains(t, upSQL, "sepay_id",
		"transactions table should have sepay_id column")
	assert.Regexp(t, regexp.MustCompile(`sepay_id\s+VARCHAR\(`), upSQL,
		"sepay_id should be VARCHAR type")

	// Test: transactions table has reference_code column
	assert.Contains(t, upSQL, "reference_code",
		"transactions table should have reference_code column")

	// Test: transactions table has webhook_payload column
	assert.Contains(t, upSQL, "webhook_payload",
		"transactions table should have webhook_payload column")
	assert.Regexp(t, regexp.MustCompile(`webhook_payload\s+JSONB`), upSQL,
		"webhook_payload should be JSONB type")

	// Test: Contains CREATE TABLE user_series_purchases statement
	assert.Regexp(t, regexp.MustCompile(`CREATE TABLE( IF NOT EXISTS)? user_series_purchases`), upSQL,
		"Up migration should contain CREATE TABLE user_series_purchases")

	// Test: user_series_purchases has composite primary key
	assert.Regexp(t, regexp.MustCompile(`PRIMARY KEY\s*\(\s*user_id\s*,\s*series_id\s*\)`), upSQL,
		"user_series_purchases should have composite primary key (user_id, series_id)")

	// Test: Contains ALTER TABLE subscriptions
	assert.Regexp(t, regexp.MustCompile(`ALTER TABLE subscriptions`), upSQL,
		"Up migration should contain ALTER TABLE subscriptions")

	// Test: subscriptions table gets expires_at column (accounting for IF NOT EXISTS)
	assert.Regexp(t, regexp.MustCompile(`ADD COLUMN( IF NOT EXISTS)?\s+expires_at\s+TIMESTAMP`), upSQL,
		"ALTER TABLE should add expires_at TIMESTAMP column")

	// Test: subscriptions table gets tier column (accounting for IF NOT EXISTS)
	assert.Regexp(t, regexp.MustCompile(`ADD COLUMN( IF NOT EXISTS)?\s+tier\s+VARCHAR`), upSQL,
		"ALTER TABLE should add tier VARCHAR column")

	// Test: subscriptions table gets updated_at column (accounting for IF NOT EXISTS)
	assert.Regexp(t, regexp.MustCompile(`ADD COLUMN( IF NOT EXISTS)?\s+updated_at\s+TIMESTAMP`), upSQL,
		"ALTER TABLE should add updated_at TIMESTAMP column")
}

// TestUpMigrationHasRequiredIndexes verifies the up migration creates required indexes
func TestUpMigrationHasRequiredIndexes(t *testing.T) {
	upSQL := readUpMigration(t)

	// Test: Index on transactions.user_id
	assert.Contains(t, upSQL, "ON transactions(user_id)",
		"Should create index on transactions.user_id")

	// Test: Index on transactions.status
	assert.Contains(t, upSQL, "ON transactions(status)",
		"Should create index on transactions.status")

	// Test: Unique index on transactions.sepay_id
	assert.Contains(t, upSQL, "UNIQUE INDEX IF NOT EXISTS idx_transactions_sepay_id ON transactions(sepay_id)",
		"Should create unique index on transactions.sepay_id")

	// Test: Index on transactions.reference_code
	assert.Contains(t, upSQL, "ON transactions(reference_code)",
		"Should create index on transactions.reference_code")

	// Test: Index on transactions.order_id
	assert.Contains(t, upSQL, "ON transactions(order_id)",
		"Should create index on transactions.order_id")

	// Test: Index on subscriptions.expires_at
	assert.Contains(t, upSQL, "ON subscriptions(expires_at)",
		"Should create index on subscriptions.expires_at")
}

// TestUpMigrationHasTransactionsTableFields verifies transactions table has all required fields
func TestUpMigrationHasTransactionsTableFields(t *testing.T) {
	upSQL := strings.ToLower(readUpMigration(t))

	requiredFields := []string{
		"id uuid primary key",
		"user_id uuid not null",
		"amount decimal",
		"currency varchar",
		"provider varchar",
		"gateway varchar",
		"type varchar",
		"status varchar",
		"target_id uuid",
		"plan_id varchar",
		"content text",
		"order_id varchar",
		"created_at timestamp",
		"updated_at timestamp",
	}

	for _, field := range requiredFields {
		assert.Contains(t, upSQL, field,
			"transactions table should have field: %s", field)
	}
}

// TestUpMigrationHasUserSeriesPurchasesTableFields verifies user_series_purchases table has all required fields
func TestUpMigrationHasUserSeriesPurchasesTableFields(t *testing.T) {
	upSQL := strings.ToLower(readUpMigration(t))

	requiredFields := []string{
		"user_id uuid",
		"series_id uuid",
		"amount decimal",
		"created_at timestamp",
	}

	for _, field := range requiredFields {
		assert.Contains(t, upSQL, field,
			"user_series_purchases table should have field: %s", field)
	}
}

// TestDownMigrationContainsRollbackStatements verifies the down migration has proper rollback
func TestDownMigrationContainsRollbackStatements(t *testing.T) {
	downSQL := strings.ToLower(readDownMigration(t))

	// Test: Drops index on subscriptions.expires_at
	assert.Contains(t, downSQL, "idx_subscriptions_expires_at",
		"Down migration should drop index idx_subscriptions_expires_at")

	// Test: Drops transactions table
	assert.Contains(t, downSQL, "drop table if exists transactions",
		"Down migration should drop transactions table")

	// Test: Drops user_series_purchases table
	assert.Contains(t, downSQL, "drop table if exists user_series_purchases",
		"Down migration should drop user_series_purchases table")
}

// TestDownMigrationAltersSubscriptions verifies the down migration reverts subscriptions changes
func TestDownMigrationAltersSubscriptions(t *testing.T) {
	downSQL := strings.ToLower(readDownMigration(t))

	// Test: Drops added columns from subscriptions
	assert.Contains(t, downSQL, "drop column",
		"Down migration should DROP columns from subscriptions")

	// Verify columns to drop
	assert.Contains(t, downSQL, "expires_at",
		"Down migration should drop expires_at column")
	assert.Contains(t, downSQL, "tier",
		"Down migration should drop tier column")
	assert.Contains(t, downSQL, "updated_at",
		"Down migration should drop updated_at column")
}

// TestMigrationFileNamesFollowConvention verifies migration files follow naming convention
func TestMigrationFileNamesFollowConvention(t *testing.T) {
	migrationsDir := MIGRATIONS_DIR
	entries, err := os.ReadDir(migrationsDir)
	require.NoError(t, err, "Should be able to read migrations directory")

	// Check for expected up file
	foundUp := false
	foundDown := false

	for _, entry := range entries {
		if entry.Name() == UP_FILE_PATTERN {
			foundUp = true
		}
		if entry.Name() == DOWN_FILE_PATTERN {
			foundDown = true
		}
	}

	assert.True(t, foundUp, "Should find up migration file: %s", UP_FILE_PATTERN)
	assert.True(t, foundDown, "Should find down migration file: %s", DOWN_FILE_PATTERN)
}
