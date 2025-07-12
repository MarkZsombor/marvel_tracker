package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Marvel Champions Play Tracker",
	})
}

func Plays(c *gin.Context) {
	c.HTML(http.StatusOK, "plays.html", gin.H{
		"title": "Plays",
		"plays": []gin.H{},
	})
}

func NewPlay(c *gin.Context) {
	c.HTML(http.StatusOK, "new_play.html", gin.H{
		"title": "New Play",
	})
}
