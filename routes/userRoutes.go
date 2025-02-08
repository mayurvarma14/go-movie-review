package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/controllers"
	"github.com/mayurvarma14/go-movie-review/middleware"
)

func UserRoutes(router *gin.Engine, uc *controllers.UserController) {
	router.Use(middleware.AuthenticateUser())
	router.GET("/users/:user_id", uc.GetUser()) // Get a specific user
	router.GET("/users", uc.GetUsers())         // Get all users (admin only)
}
