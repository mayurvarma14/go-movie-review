package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/controllers"
	"github.com/mayurvarma14/go-movie-review/middleware"
)

func GenreRoutes(router *gin.Engine, gc *controllers.GenreController) {
	router.Use(middleware.AuthenticateUser())
	router.POST("/genres", gc.CreateGenre())
	router.GET("/genres/:genre_id", gc.GetGenre())
	router.GET("/genres", gc.GetGenres())
	router.PUT("/genres/:genre_id", gc.EditGenre())
	router.DELETE("/genres/:genre_id", gc.DeleteGenre())

}
