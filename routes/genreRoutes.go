package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/controllers"
	"github.com/mayurvarma14/go-movie-review/middleware"
)

func GenreRoutes(router *gin.Engine, gc *controllers.GenreController) {
	router.Use(middleware.AuthenticateUser())
	router.POST("/genres", gc.CreateGenre())             // Create a new genre (admin only)
	router.GET("/genres/:genre_id", gc.GetGenre())       // Get a specific genre
	router.GET("/genres", gc.GetGenres())                // Get all genres
	router.PUT("/genres/:genre_id", gc.EditGenre())      // Update a genre (admin only)
	router.DELETE("/genres/:genre_id", gc.DeleteGenre()) // Delete a genre (admin only)
}
