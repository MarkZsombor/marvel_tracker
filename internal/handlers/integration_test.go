package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupIntegrationTestRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	gin.SetMode(gin.TestMode)

	// Create in-memory database for integration tests
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

	// Insert test data
	_, err = db.Exec("INSERT INTO scenarios (id, name) VALUES (1, 'Rhino'), (2, 'Klaw')")
	require.NoError(t, err)

	r := gin.New()

	// Load templates for integration testing
	r.LoadHTMLFiles(
		"../../templates/index.html",
		"../../templates/plays.html",
		"../../templates/new_play.html",
		"../../templates/error.html",
	)

	return r, db
}

func TestHandlers_Integration(t *testing.T) {
	r, db := setupIntegrationTestRouter(t)
	defer db.Close()

	// Setup routes
	r.GET("/", Home)
	r.GET("/plays", Plays)
	r.GET("/plays/new", NewPlay)

	t.Run("Full Navigation Flow", func(t *testing.T) {
		// Test home page
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Welcome to Marvel Champions Play Tracker")
		assert.Contains(t, w.Body.String(), "View Plays")
		assert.Contains(t, w.Body.String(), "Log New Play")

		// Test plays page
		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, "/plays", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Play History")
		assert.Contains(t, w.Body.String(), "No plays recorded yet")

		// Test new play page
		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, "/plays/new", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Log New Play")
		assert.Contains(t, w.Body.String(), "Date")
		assert.Contains(t, w.Body.String(), "Scenario")
		assert.Contains(t, w.Body.String(), "Difficulty")
		assert.Contains(t, w.Body.String(), "Outcome")
	})

	t.Run("Template Variables", func(t *testing.T) {
		// Test that title variables are properly passed to templates
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "<title>Marvel Champions Play Tracker - Marvel Champions Play Tracker</title>")

		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, "/plays", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "<title>Plays - Marvel Champions Play Tracker</title>")

		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, "/plays/new", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "<title>New Play - Marvel Champions Play Tracker</title>")
	})

	t.Run("Navigation Links Present", func(t *testing.T) {
		pages := []string{"/", "/plays", "/plays/new"}

		for _, page := range pages {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, page, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			// Verify navigation links are present on all pages
			body := w.Body.String()
			assert.Contains(t, body, `href="/"`)
			assert.Contains(t, body, `href="/plays"`)
			assert.Contains(t, body, `href="/plays/new"`)
			assert.Contains(t, body, "Marvel Champions Play Tracker")
		}
	})

	t.Run("Form Elements on New Play Page", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/plays/new", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()

		// Verify form exists with correct action and method
		assert.Contains(t, body, `action="/plays"`)
		assert.Contains(t, body, `method="POST"`)

		// Verify form fields
		assert.Contains(t, body, `name="date"`)
		assert.Contains(t, body, `name="scenario"`)
		assert.Contains(t, body, `name="difficulty"`)
		assert.Contains(t, body, `name="outcome"`)
		assert.Contains(t, body, `name="notes"`)

		// Verify difficulty options
		assert.Contains(t, body, "Standard I")
		assert.Contains(t, body, "Standard II")
		assert.Contains(t, body, "Expert I")
		assert.Contains(t, body, "Expert II")
		assert.Contains(t, body, "Heroic I")

		// Verify outcome options
		assert.Contains(t, body, `value="win"`)
		assert.Contains(t, body, `value="loss"`)

		// Verify buttons
		assert.Contains(t, body, "Save Play")
		assert.Contains(t, body, "Cancel")
	})

	t.Run("HTMX and Tailwind Resources", func(t *testing.T) {
		pages := []string{"/", "/plays", "/plays/new"}

		for _, page := range pages {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, page, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			body := w.Body.String()

			// Verify HTMX is loaded
			assert.Contains(t, body, "https://unpkg.com/htmx.org")

			// Verify Tailwind CSS is loaded
			assert.Contains(t, body, "https://cdn.tailwindcss.com")
		}
	})
}

func TestHandlers_ErrorScenarios(t *testing.T) {
	r, db := setupIntegrationTestRouter(t)
	defer db.Close()

	r.GET("/", Home)
	r.GET("/plays", Plays)
	r.GET("/plays/new", NewPlay)

	t.Run("Non-existent Route", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/nonexistent", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Invalid HTTP Method", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/", nil)
		r.ServeHTTP(w, req)

		// Gin returns 404 for routes that don't match method, not 405
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandlers_ResponseHeaders(t *testing.T) {
	r, db := setupIntegrationTestRouter(t)
	defer db.Close()

	r.GET("/", Home)

	t.Run("Content Type Headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
	})
}
