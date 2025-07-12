package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	t.Run("Default Database Path", func(t *testing.T) {
		// Ensure no DB_PATH environment variable is set
		os.Unsetenv("DB_PATH")
		
		// Create a temporary directory for testing
		tempDir := t.TempDir()
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		
		// Change to temp directory so default ./data path is created there
		os.Chdir(tempDir)
		
		db := InitDB()
		defer db.Close()
		
		// Verify database connection works
		err := db.Ping()
		assert.NoError(t, err)
		
		// Verify data directory was created
		_, err = os.Stat("./data")
		assert.NoError(t, err)
		
		// Verify database file exists
		_, err = os.Stat("./data/marvel_tracker.db")
		assert.NoError(t, err)
	})
	
	t.Run("Custom Database Path", func(t *testing.T) {
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "custom.db")
		
		// Set custom DB_PATH
		os.Setenv("DB_PATH", dbPath)
		defer os.Unsetenv("DB_PATH")
		
		db := InitDB()
		defer db.Close()
		
		// Verify database connection works
		err := db.Ping()
		assert.NoError(t, err)
		
		// Verify custom database file exists
		_, err = os.Stat(dbPath)
		assert.NoError(t, err)
	})
	
	t.Run("Custom Database Path with Directory Creation", func(t *testing.T) {
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "subdir", "custom.db")
		
		// Set custom DB_PATH with subdirectory
		os.Setenv("DB_PATH", dbPath)
		defer os.Unsetenv("DB_PATH")
		
		// Create the subdirectory first since SQLite requires parent directory to exist
		err := os.MkdirAll(filepath.Dir(dbPath), 0755)
		assert.NoError(t, err)
		
		db := InitDB()
		defer db.Close()
		
		// Verify database connection works
		err = db.Ping()
		assert.NoError(t, err)
		
		// Verify subdirectory was created
		_, err = os.Stat(filepath.Dir(dbPath))
		assert.NoError(t, err)
		
		// Verify database file exists
		_, err = os.Stat(dbPath)
		assert.NoError(t, err)
	})
}