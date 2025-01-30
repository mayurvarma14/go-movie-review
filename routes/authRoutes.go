package routes

import (
	controllers "github.com/mayurvarma14/go-movie-review/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, uc *controllers.UserController) {
	router.POST("users/signup", uc.SignUp())
	router.POST("users/login", uc.Login())
}
