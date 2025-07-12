package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestRouter initializes a Gin router for testing, loading all templates.
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Load templates explicitly to avoid conflicts between content blocks
	r.LoadHTMLFiles(
		"../../templates/index.html",
		"../../templates/plays.html", 
		"../../templates/new_play.html",
		"../../templates/error.html",
	)

	return r
}

func TestHomeHandler(t *testing.T) {
	r := setupTestRouter()
	r.GET("/", Home)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	t.Logf("Response Body for Home Handler:\n%s", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Welcome to Marvel Champions Play Tracker")
}

func TestPlaysHandler(t *testing.T) {
	r := setupTestRouter()
	r.GET("/plays", Plays)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/plays", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Play History")
	assert.Contains(t, w.Body.String(), "No plays recorded yet.")
}

func TestNewPlayHandler(t *testing.T) {
	r := setupTestRouter()
	r.GET("/plays/new", NewPlay)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/plays/new", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Log New Play")
}
