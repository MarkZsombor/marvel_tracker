package config

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	return db
}

func createTestMigrationFiles(t *testing.T, dir string) {
	// Create migration directory
	err := os.MkdirAll(dir, 0755)
	require.NoError(t, err)

	// Create test migration files
	migrations := map[string]string{
		"001_create_users.sql": `
-- Create users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);
`,
		"002_create_posts.sql": `
-- Create posts table
CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
`,
		"003_add_email_to_users.sql": `
-- Add email column to users
ALTER TABLE users ADD COLUMN email TEXT;
`,
	}

	for filename, content := range migrations {
		filePath := filepath.Join(dir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)
	}
}

func TestCreateMigrationsTable(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := createMigrationsTable(db)
	assert.NoError(t, err)

	// Verify table was created
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='migrations'").Scan(&tableName)
	assert.NoError(t, err)
	assert.Equal(t, "migrations", tableName)

	// Verify table structure
	rows, err := db.Query("PRAGMA table_info(migrations)")
	assert.NoError(t, err)
	defer rows.Close()

	columns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var defaultValue sql.NullString

		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		require.NoError(t, err)
		columns[name] = true
	}

	assert.True(t, columns["id"])
	assert.True(t, columns["filename"])
	assert.True(t, columns["applied_at"])
}

func TestRunMigrations(t *testing.T) {
	t.Run("Successful Migration", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// Create temporary migration directory
		tempDir := t.TempDir()
		migrationDir := filepath.Join(tempDir, "migrations")
		createTestMigrationFiles(t, migrationDir)

		// Change working directory temporarily
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tempDir)

		err := RunMigrations(db)
		assert.NoError(t, err)

		// Verify migrations table exists
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 3, count) // Should have 3 migration records

		// Verify actual tables were created
		tables := []string{"users", "posts"}
		for _, table := range tables {
			var tableName string
			err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&tableName)
			assert.NoError(t, err)
			assert.Equal(t, table, tableName)
		}

		// Verify email column was added to users
		var hasEmail bool
		rows, err := db.Query("PRAGMA table_info(users)")
		assert.NoError(t, err)
		defer rows.Close()

		for rows.Next() {
			var cid int
			var name, dataType string
			var notNull, pk int
			var defaultValue sql.NullString

			err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
			require.NoError(t, err)
			if name == "email" {
				hasEmail = true
				break
			}
		}
		assert.True(t, hasEmail)
	})

	t.Run("No Migration Files", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// Create empty migration directory
		tempDir := t.TempDir()
		migrationDir := filepath.Join(tempDir, "migrations")
		err := os.MkdirAll(migrationDir, 0755)
		require.NoError(t, err)

		// Change working directory temporarily
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tempDir)

		err = RunMigrations(db)
		assert.NoError(t, err)

		// Should still create migrations table
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Skip Already Applied Migrations", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// Create temporary migration directory
		tempDir := t.TempDir()
		migrationDir := filepath.Join(tempDir, "migrations")
		createTestMigrationFiles(t, migrationDir)

		// Change working directory temporarily
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tempDir)

		// Run migrations first time
		err := RunMigrations(db)
		assert.NoError(t, err)

		// Run migrations second time - should skip all
		err = RunMigrations(db)
		assert.NoError(t, err)

		// Should still only have 3 migration records
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("Invalid SQL in Migration", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// Create temporary migration directory with invalid SQL
		tempDir := t.TempDir()
		migrationDir := filepath.Join(tempDir, "migrations")
		err := os.MkdirAll(migrationDir, 0755)
		require.NoError(t, err)

		// Create migration with invalid SQL
		invalidMigration := `
-- Invalid SQL
CREATE INVALID TABLE syntax;
`
		filePath := filepath.Join(migrationDir, "001_invalid.sql")
		err = os.WriteFile(filePath, []byte(invalidMigration), 0644)
		require.NoError(t, err)

		// Change working directory temporarily
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tempDir)

		err = RunMigrations(db)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "syntax error")
	})

	t.Run("Migration with Comments and Empty Lines", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// Create temporary migration directory
		tempDir := t.TempDir()
		migrationDir := filepath.Join(tempDir, "migrations")
		err := os.MkdirAll(migrationDir, 0755)
		require.NoError(t, err)

		// Create migration with comments and empty lines
		migrationWithComments := `
-- This is a comment
-- Another comment

CREATE TABLE test_table (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Inline comment
    name TEXT NOT NULL
);

-- Final comment

`
		filePath := filepath.Join(migrationDir, "001_with_comments.sql")
		err = os.WriteFile(filePath, []byte(migrationWithComments), 0644)
		require.NoError(t, err)

		// Change working directory temporarily
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tempDir)

		err = RunMigrations(db)
		assert.NoError(t, err)

		// Verify table was created despite comments
		var tableName string
		err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&tableName)
		assert.NoError(t, err)
		assert.Equal(t, "test_table", tableName)
	})

	t.Run("No Migrations Directory", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// Use a directory without migrations folder
		tempDir := t.TempDir()
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tempDir)

		err := RunMigrations(db)
		assert.NoError(t, err) // Should not error, just find no files

		// Should still create migrations table
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}
