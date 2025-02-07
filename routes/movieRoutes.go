package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/controllers"
	"github.com/mayurvarma14/go-movie-review/middleware"
)

func MovieRoutes(router *gin.Engine, mc *controllers.MovieController) {
	router.Use(middleware.AuthenticateUser())
	router.POST("/movies", mc.CreateMovie())
	router.GET("/movies/:movie_id", mc.GetMovie())
	router.GET("/movies", mc.GetMovies())
	router.PUT("/movies/:movie_id", mc.UpdateMovie())
	router.GET("/movies/search", mc.SearchMovieByQuery())
	router.GET("movies/filter", mc.SearchMovieByGenre())
	router.DELETE("/movies/:movie_id", mc.DeleteMovie())
}
