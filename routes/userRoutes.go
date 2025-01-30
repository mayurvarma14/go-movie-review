package routes

import (
	"github.com/mayurvarma14/go-movie-review/controllers"
	"github.com/mayurvarma14/go-movie-review/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, uc *controllers.UserController) {
	router.Use(middleware.AuthenticateUser())
	router.GET("/users/:user_id", uc.GetUser())
	router.GET("/users", uc.GetUsers())
}
