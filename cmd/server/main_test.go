package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"marvel_tracker/internal/config"
	"marvel_tracker/internal/handlers"
)

func setupTestServer(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Create test database
	tempDir := t.TempDir()
	dbPath := tempDir + "/test.db"
	os.Setenv("DB_PATH", dbPath)
	defer os.Unsetenv("DB_PATH")

	db := config.InitDB()
	defer db.Close()

	err := config.RunMigrations(db)
	require.NoError(t, err)

	r := gin.New()

	// Load templates (create minimal test templates)
	testTemplatesDir := tempDir + "/templates"
	err = os.MkdirAll(testTemplatesDir, 0755)
	require.NoError(t, err)

	// Create minimal test templates
	templates := map[string]string{
		"index.html":    `<!DOCTYPE html><html><head><title>{{.title}}</title></head><body><h1>Home</h1></body></html>`,
		"plays.html":    `<!DOCTYPE html><html><head><title>{{.title}}</title></head><body><h1>Plays</h1></body></html>`,
		"new_play.html": `<!DOCTYPE html><html><head><title>{{.title}}</title></head><body><h1>New Play</h1></body></html>`,
		"error.html":    `<!DOCTYPE html><html><head><title>{{.title}}</title></head><body><h1>Error</h1></body></html>`,
	}

	for filename, content := range templates {
		err = os.WriteFile(testTemplatesDir+"/"+filename, []byte(content), 0644)
		require.NoError(t, err)
	}

	r.LoadHTMLGlob(testTemplatesDir + "/*")

	// Setup routes like in main
	r.GET("/", handlers.Home)
	r.GET("/plays", handlers.Plays)
	r.GET("/plays/new", handlers.NewPlay)

	return r
}

func TestServerInitialization(t *testing.T) {
	t.Run("Server Routes Registration", func(t *testing.T) {
		server := setupTestServer(t)

		// Test that all expected routes are registered and working
		routes := []struct {
			method         string
			path           string
			expectedStatus int
		}{
			{"GET", "/", http.StatusOK},
			{"GET", "/plays", http.StatusOK},
			{"GET", "/plays/new", http.StatusOK},
		}

		for _, route := range routes {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(route.method, route.path, nil)
			server.ServeHTTP(w, req)

			assert.Equal(t, route.expectedStatus, w.Code,
				"Route %s %s should return status %d", route.method, route.path, route.expectedStatus)
		}
	})

	t.Run("Template Loading", func(t *testing.T) {
		server := setupTestServer(t)

		// Test that templates are loaded and rendering correctly
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		server.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "<title>Marvel Champions Play Tracker</title>")
		assert.Contains(t, w.Body.String(), "<h1>Home</h1>")
	})

	t.Run("Database Integration", func(t *testing.T) {
		// Test that database initialization and migrations work
		tempDir := t.TempDir()
		dbPath := tempDir + "/integration_test.db"

		os.Setenv("DB_PATH", dbPath)
		defer os.Unsetenv("DB_PATH")

		db := config.InitDB()
		defer db.Close()

		// Test database connection
		err := db.Ping()
		assert.NoError(t, err)

		// Test migrations (skip table verification since migrations directory may not exist in test)
		err = config.RunMigrations(db)
		assert.NoError(t, err)

		// Verify migrations table exists (this should always be created)
		var tableName string
		err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='migrations'").Scan(&tableName)
		assert.NoError(t, err, "Migrations table should exist")
		assert.Equal(t, "migrations", tableName)
	})

	t.Run("Static File Serving Configuration", func(t *testing.T) {
		server := setupTestServer(t)

		// Create a test static directory and file
		tempDir := t.TempDir()
		staticDir := tempDir + "/static"
		err := os.MkdirAll(staticDir, 0755)
		require.NoError(t, err)

		testCSS := "body { background: red; }"
		err = os.WriteFile(staticDir+"/test.css", []byte(testCSS), 0644)
		require.NoError(t, err)

		// Add static file serving to test server
		server.Static("/static", staticDir)

		// Test static file access
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/static/test.css", nil)
		server.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, testCSS, w.Body.String())
		assert.Contains(t, w.Header().Get("Content-Type"), "text/css")
	})

	t.Run("Error Handling in Production", func(t *testing.T) {
		server := setupTestServer(t)

		// Test 404 handling
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/nonexistent", nil)
		server.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Environment Variable Configuration", func(t *testing.T) {
		// Test different DB_PATH configurations
		testCases := []struct {
			name   string
			dbPath string
		}{
			{"Custom Path", "/tmp/custom_test.db"},
			{"Relative Path", "./custom_test.db"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Clean up any existing file
				os.Remove(tc.dbPath)
				defer os.Remove(tc.dbPath)

				os.Setenv("DB_PATH", tc.dbPath)
				defer os.Unsetenv("DB_PATH")

				db := config.InitDB()
				defer db.Close()

				err := db.Ping()
				assert.NoError(t, err)

				// Verify file was created at the expected path
				_, err = os.Stat(tc.dbPath)
				assert.NoError(t, err)
			})
		}
	})
}

func TestServerPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	t.Run("Response Time Benchmarks", func(t *testing.T) {
		server := setupTestServer(t)

		routes := []string{"/", "/plays", "/plays/new"}

		for _, route := range routes {
			start := time.Now()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", route, nil)
			server.ServeHTTP(w, req)

			duration := time.Since(start)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Less(t, duration, 100*time.Millisecond,
				"Route %s should respond within 100ms, took %v", route, duration)
		}
	})

	t.Run("Concurrent Request Handling", func(t *testing.T) {
		server := setupTestServer(t)

		const numRequests = 10
		results := make(chan bool, numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/", nil)
				server.ServeHTTP(w, req)
				results <- w.Code == http.StatusOK
			}()
		}

		// Collect results
		successCount := 0
		for i := 0; i < numRequests; i++ {
			if <-results {
				successCount++
			}
		}

		assert.Equal(t, numRequests, successCount,
			"All concurrent requests should succeed")
	})
}
