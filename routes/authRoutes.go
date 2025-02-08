package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/controllers"
)

func AuthRoutes(router *gin.Engine, uc *controllers.UserController) {
	router.POST("/users/signup", uc.SignUp())
	router.POST("/users/login", uc.Login())
}
