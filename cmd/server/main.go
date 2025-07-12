package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"marvel_tracker/internal/config"
	"marvel_tracker/internal/handlers"
)

func main() {
	db := config.InitDB()
	defer db.Close()

	if err := config.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", handlers.Home)
	r.GET("/plays", handlers.Plays)
	r.GET("/plays/new", handlers.NewPlay)

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}