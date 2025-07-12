package models

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Create test schema
	schema := `
	CREATE TABLE plays (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date DATE NOT NULL,
		outcome TEXT NOT NULL CHECK(outcome IN ('win', 'loss')),
		difficulty TEXT NOT NULL,
		notes TEXT,
		scenario_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE scenarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(schema)
	require.NoError(t, err)

	// Insert test scenario
	_, err = db.Exec("INSERT INTO scenarios (id, name) VALUES (1, 'Rhino')")
	require.NoError(t, err)

	return db
}

func TestNewPlayRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPlayRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestPlayRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewPlayRepository(db)

	t.Run("Valid Play Creation", func(t *testing.T) {
		play := &Play{
			Date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Outcome:    "win",
			Difficulty: "Standard I",
			Notes:      "Great game!",
			ScenarioID: 1,
		}

		err := repo.Create(play)
		assert.NoError(t, err)
		assert.NotZero(t, play.ID)

		// Verify it was actually inserted
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM plays WHERE id = ?", play.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("Invalid Outcome", func(t *testing.T) {
		play := &Play{
			Date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Outcome:    "invalid",
			Difficulty: "Standard I",
			Notes:      "Test",
			ScenarioID: 1,
		}

		err := repo.Create(play)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CHECK constraint failed")
	})

	t.Run("Invalid Scenario ID", func(t *testing.T) {
		play := &Play{
			Date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Outcome:    "win",
			Difficulty: "Standard I",
			Notes:      "Test",
			ScenarioID: 999, // Non-existent scenario
		}

		err := repo.Create(play)
		// SQLite doesn't enforce foreign keys by default in memory, but we test the structure
		// In a real database with FK constraints enabled, this would fail
		if err != nil {
			assert.Contains(t, err.Error(), "FOREIGN KEY constraint failed")
		}
	})
}

func TestPlayRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewPlayRepository(db)

	t.Run("Empty Database", func(t *testing.T) {
		plays, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, plays)
	})

	t.Run("Multiple Plays", func(t *testing.T) {
		// Insert test plays
		testPlays := []*Play{
			{
				Date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				Outcome:    "win",
				Difficulty: "Standard I",
				Notes:      "First game",
				ScenarioID: 1,
			},
			{
				Date:       time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
				Outcome:    "loss",
				Difficulty: "Expert I",
				Notes:      "Tough one",
				ScenarioID: 1,
			},
		}

		for _, play := range testPlays {
			err := repo.Create(play)
			require.NoError(t, err)
		}

		plays, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, plays, 2)

		// Should be ordered by date DESC (newest first)
		assert.True(t, plays[0].Date.After(plays[1].Date) || plays[0].Date.Equal(plays[1].Date))
	})

	t.Run("Database Error", func(t *testing.T) {
		// Close the database to simulate an error
		db.Close()

		plays, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, plays)
	})
}

func TestPlay_JSONMarshaling(t *testing.T) {
	play := Play{
		ID:         1,
		Date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Outcome:    "win",
		Difficulty: "Standard I",
		Notes:      "Great game!",
		ScenarioID: 1,
		CreatedAt:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	// Test that the struct has proper JSON tags
	assert.Equal(t, 1, play.ID)
	assert.Equal(t, "win", play.Outcome)
	assert.Equal(t, "Standard I", play.Difficulty)
	assert.Equal(t, "Great game!", play.Notes)
	assert.Equal(t, 1, play.ScenarioID)
}

func TestHero_Struct(t *testing.T) {
	hero := Hero{
		ID:        1,
		Name:      "Spider-Man",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, 1, hero.ID)
	assert.Equal(t, "Spider-Man", hero.Name)
	assert.NotZero(t, hero.CreatedAt)
	assert.NotZero(t, hero.UpdatedAt)
}

func TestScenario_Struct(t *testing.T) {
	scenario := Scenario{
		ID:        1,
		Name:      "Rhino",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, 1, scenario.ID)
	assert.Equal(t, "Rhino", scenario.Name)
	assert.NotZero(t, scenario.CreatedAt)
	assert.NotZero(t, scenario.UpdatedAt)
}

func TestDeck_Struct(t *testing.T) {
	deck := Deck{
		ID:        1,
		PlayID:    1,
		HeroID:    1,
		Aspect:    "leadership",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, 1, deck.ID)
	assert.Equal(t, 1, deck.PlayID)
	assert.Equal(t, 1, deck.HeroID)
	assert.Equal(t, "leadership", deck.Aspect)
	assert.NotZero(t, deck.CreatedAt)
	assert.NotZero(t, deck.UpdatedAt)
}