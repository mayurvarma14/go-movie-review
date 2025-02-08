package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/controllers"
	"github.com/mayurvarma14/go-movie-review/middleware"
)

func MovieRoutes(router *gin.Engine, mc *controllers.MovieController) {
	router.Use(middleware.AuthenticateUser())
	router.POST("/movies", mc.CreateMovie())              // Create a new movie (admin only)
	router.GET("/movies/:movie_id", mc.GetMovie())        // Get a specific movie
	router.GET("/movies", mc.GetMovies())                 // Get all movies
	router.PUT("/movies/:movie_id", mc.UpdateMovie())     // Update a movie (admin only)
	router.GET("/movies/search", mc.SearchMovieByQuery()) // Search movies by name
	router.GET("/movies/filter", mc.SearchMovieByGenre()) // Search movies by genre
	router.DELETE("/movies/:movie_id", mc.DeleteMovie())  // Delete a movie (admin only)
}
