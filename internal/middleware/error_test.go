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

	t.Run("Multiple Errors in Chain", func(t *testing.T) {
		r := setupRouter()
		r.GET("/", func(c *gin.Context) {
			c.Error(errors.New("first error"))
			c.Error(errors.New("second error"))
			c.Error(errors.New("last error"))
			c.AbortWithStatus(http.StatusInternalServerError)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal Server Error")
		// The middleware should handle the last error in the chain
	})

	t.Run("Error Without Status Set", func(t *testing.T) {
		r := setupRouter()
		r.GET("/", func(c *gin.Context) {
			c.Error(errors.New("error without status"))
			// Don't set any status - should remain 200
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		// When no status is set but errors exist, response should still be 200
		// The middleware only acts on errors when status codes indicate errors
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Error Template Data", func(t *testing.T) {
		r := setupRouter()
		r.GET("/", func(c *gin.Context) {
			c.Error(errors.New("test error"))
			c.AbortWithStatus(http.StatusNotFound)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		body := w.Body.String()
		
		// Verify template variables are correctly populated
		assert.Contains(t, body, "Page Not Found") // title
		assert.Contains(t, body, "The page you&#39;re looking for doesn&#39;t exist.") // message (HTML escaped)
		assert.Contains(t, body, "404") // code
		assert.Contains(t, body, "Back to Home") // link text
	})

	t.Run("Different Error Status Codes", func(t *testing.T) {
		testCases := []struct {
			status       int
			expectedText string
		}{
			{http.StatusUnauthorized, "An Unexpected Error Occurred"},
			{http.StatusForbidden, "An Unexpected Error Occurred"},
			{http.StatusBadRequest, "An Unexpected Error Occurred"},
			{http.StatusConflict, "An Unexpected Error Occurred"},
			{http.StatusServiceUnavailable, "An Unexpected Error Occurred"},
		}

		for _, tc := range testCases {
			r := setupRouter()
			r.GET("/", func(c *gin.Context) {
				c.Error(errors.New("test error"))
				c.AbortWithStatus(tc.status)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.status, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedText)
		}
	})

	t.Run("Success Status with Error Should Not Trigger Handler", func(t *testing.T) {
		r := setupRouter()
		r.GET("/", func(c *gin.Context) {
			c.Error(errors.New("error with success status"))
			c.String(http.StatusOK, "Success Response")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		// Should return the success response, error handler still processes but doesn't override
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success Response")
	})

	t.Run("Middleware Chain Order", func(t *testing.T) {
		r := setupRouter()
		
		middlewareCalled := false
		
		// Add another middleware after error handler
		r.Use(func(c *gin.Context) {
			middlewareCalled = true
			c.Next()
		})
		
		r.GET("/", func(c *gin.Context) {
			c.Error(errors.New("test error"))
			c.AbortWithStatus(http.StatusInternalServerError)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.True(t, middlewareCalled)
		assert.Contains(t, w.Body.String(), "Internal Server Error")
	})
}
