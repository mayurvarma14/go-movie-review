package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/controllers"
	"github.com/mayurvarma14/go-movie-review/middleware"
)

func ReviewRoutes(router *gin.Engine, rc *controllers.ReviewController) {
	router.Use(middleware.AuthenticateUser())
	router.POST("reviews", rc.AddReview())
	router.GET("reviews/filter", rc.ViewAMovieReviews())
	router.DELETE("reviews/:id", rc.DeleteReview())
	router.GET("reviews/user/:reviewer_id", rc.AllUserReviews())
}
