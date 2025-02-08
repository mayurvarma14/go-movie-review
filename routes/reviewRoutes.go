package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/controllers"
	"github.com/mayurvarma14/go-movie-review/middleware"
)

func ReviewRoutes(router *gin.Engine, rc *controllers.ReviewController) {
	router.Use(middleware.AuthenticateUser())
	router.POST("/reviews", rc.AddReview())                       // Add a review (user only)
	router.GET("/reviews/filter", rc.ViewAMovieReviews())         // Get reviews for a movie
	router.DELETE("/reviews/:id", rc.DeleteReview())              // Delete a review (owner or admin)
	router.GET("/reviews/user/:reviewer_id", rc.AllUserReviews()) // Get all reviews by a user
}
