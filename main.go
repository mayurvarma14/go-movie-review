package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/controllers"
	"github.com/mayurvarma14/go-movie-review/database"
	"github.com/mayurvarma14/go-movie-review/internals/config"
	"github.com/mayurvarma14/go-movie-review/routes"
)

func main() {
	// Load env variables
	config.LoadEnv()

	// Initialize Database
	ctx := context.Background()
	db, err := database.New(ctx)
	if err != nil {
		log.Fatal("Database init failed:", err)
	}
	defer db.Client.Disconnect(ctx)

	// Create controller with injected dependency
	uc := controllers.NewUserController(db)
	gc := controllers.NewGenreController(db)
	mc := controllers.NewMovieController(db)

	// Set up Gin
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := gin.Default()

	//Log events
	r.Use(gin.Logger())

	// Register app routes
	routes.AuthRoutes(r, uc)
	routes.UserRoutes(r, uc)
	routes.GenreRoutes(r, gc)
	routes.MovieRoutes(r, mc)

	r.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "Welcome to movie review app",
		})
	})

	r.Run(":" + port)
}
