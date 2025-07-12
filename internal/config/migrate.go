package config

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations(db *sql.DB) error {
	if err := createMigrationsTable(db); err != nil {
		return err
	}

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		filename := filepath.Base(file)

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = ?", filename).Scan(&count)
		if err != nil {
			return err
		}

		if count > 0 {
			log.Printf("Migration %s already applied, skipping", filename)
			continue
		}

		log.Printf("Running migration: %s", filename)

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		// Execute the entire migration file as one transaction
		content_str := string(content)

		// Remove comments and split by semicolon
		lines := strings.Split(content_str, "\n")
		var cleanedLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "--") {
				cleanedLines = append(cleanedLines, line)
			}
		}

		if len(cleanedLines) == 0 {
			continue
		}

		// Join lines and split by semicolon
		fullContent := strings.Join(cleanedLines, " ")
		statements := strings.Split(fullContent, ";")

		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			if _, err := db.Exec(stmt); err != nil {
				log.Printf("Error executing statement: %s", stmt)
				return err
			}
		}

		_, err = db.Exec("INSERT INTO migrations (filename) VALUES (?)", filename)
		if err != nil {
			return err
		}

		log.Printf("Migration %s completed successfully", filename)
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT NOT NULL UNIQUE,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(query)
	return err
}
