package middleware

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupRouter initializes a Gin router for testing.
func setupRouter() *gin.Engine {
	r := gin.New()

	// Correctly load all templates from the templates directory.
	// The path is relative to the project root, where `go test` is run.
	tmpl, err := template.ParseGlob("../../templates/*.html")
	if err != nil {
		panic("Failed to parse templates: " + err.Error())
	}
	r.SetHTMLTemplate(tmpl)

	// Add the error handler middleware.
	r.Use(ErrorHandler())

	return r
}

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("No Error", func(t *testing.T) {
		r := setupRouter()
		r.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Success")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Success", w.Body.String())
	})

	t.Run("Not Found Error", func(t *testing.T) {
		r := setupRouter()
		r.GET("/", func(c *gin.Context) {
			c.Error(errors.New("not found"))
			// Abort to prevent other handlers from running.
			c.AbortWithStatus(http.StatusNotFound)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Page Not Found")
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		r := setupRouter()
		r.GET("/", func(c *gin.Context) {
			c.Error(errors.New("internal server error"))
			c.AbortWithStatus(http.StatusInternalServerError)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal Server Error")
	})

	t.Run("Other Client Error", func(t *testing.T) {
		r := setupRouter()
		r.GET("/", func(c *gin.Context) {
			c.Error(errors.New("bad request"))
			c.AbortWithStatus(http.StatusBadRequest)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "An Unexpected Error Occurred")
	})
}
