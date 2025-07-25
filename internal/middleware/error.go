package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error: %v", err)

			switch c.Writer.Status() {
			case http.StatusNotFound:
				c.HTML(http.StatusNotFound, "error.html", gin.H{
					"title":   "Page Not Found",
					"message": "The page you're looking for doesn't exist.",
					"code":    404,
				})
				return
			case http.StatusInternalServerError:
				c.HTML(http.StatusInternalServerError, "error.html", gin.H{
					"title":   "Internal Server Error",
					"message": "Something went wrong on our end.",
					"code":    500,
				})
				return
			default:
				// Don't leak internal error details to the client.
				// Log the full error server-side and return a generic
				// message to the user.
				c.HTML(c.Writer.Status(), "error.html", gin.H{
					"title":   "An Unexpected Error Occurred",
					"message": "We've encountered an error and our team has been notified.",
					"code":    c.Writer.Status(),
				})
				return
			}
		}
	}
}
