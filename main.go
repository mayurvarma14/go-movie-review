package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/controllers"
	"github.com/mayurvarma14/go-movie-review/database"
	"github.com/mayurvarma14/go-movie-review/internals/config"
	"github.com/mayurvarma14/go-movie-review/routes"
)

func main() {
	config.LoadEnv()

	ctx := context.Background()
	db, err := database.New(ctx)
	if err != nil {
		log.Fatal("Database init failed:", err)
	}
	defer func() {
		if err := db.Client.Disconnect(ctx); err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
	}()

	uc := controllers.NewUserController(db)
	gc := controllers.NewGenreController(db)
	mc := controllers.NewMovieController(db)
	rc := controllers.NewReviewController(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	router.Use(gin.Logger())

	routes.AuthRoutes(router, uc)
	routes.UserRoutes(router, uc)
	routes.GenreRoutes(router, gc)
	routes.MovieRoutes(router, mc)
	routes.ReviewRoutes(router, rc)

	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the movie review API"})
	})

	log.Printf("Server listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
