package controllers

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/database"
	helper "github.com/mayurvarma14/go-movie-review/helpers"
	"github.com/mayurvarma14/go-movie-review/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReviewController struct {
	reviewCollection *mongo.Collection
}

func NewReviewController(db *database.Database) *ReviewController {
	return &ReviewController{
		reviewCollection: db.Client.Database(os.Getenv("MONGO_INITDB_DATABASE")).Collection("review"),
	}
}

func (rc *ReviewController) AddReview() gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := helper.VerifyUserType(c, "USER"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var review models.Reviews
		defer cancel()

		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationError := validate.Struct(&review); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()}})
			return
		}

		newReview := models.Reviews{
			Id:          bson.NewObjectID(),
			Movie_id:    review.Movie_id,
			Reviewer_id: c.GetString("uid"),
			Review:      review.Review,
		}

		result, err := rc.reviewCollection.InsertOne(ctx, newReview)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"Status":  http.StatusCreated,
			"Message": "success",
			"Data":    map[string]interface{}{"data": result}})
	}
}

func (rc *ReviewController) ViewAMovieReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchReviews []models.Reviews
		queryParam := c.Query("movie_id")
		if queryParam == "" {
			log.Println("No movie id passed")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		id, err := strconv.Atoi(queryParam)
		if err != nil {
			log.Println("Invalid movie id")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid movie id"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchQueryDB, err := rc.reviewCollection.Find(ctx, bson.M{"movie_id": id})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchQueryDB.All(ctx, &searchReviews)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchQueryDB.Close(ctx)
		if err := searchQueryDB.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchReviews)
	}
}

func (rc *ReviewController) DeleteReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		reviewId := c.Param("id")
		defer cancel()
		objId, _ := bson.ObjectIDFromHex(reviewId)

		result, err := rc.reviewCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				gin.H{
					" Status":  http.StatusNotFound,
					" Message": "error",
					" Data":    map[string]interface{}{"data": "Review with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			gin.H{
				"Status":  http.StatusOK,
				"Message": "success",
				"Data":    map[string]interface{}{"data": "Your review was successfully deleted!"}},
		)
	}
}

func (rc *ReviewController) AllUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchReviews []models.Reviews
		reviewId := c.Param("reviewer_id")
		if reviewId == "" {
			log.Println("No reviewer id passed")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchQueryDB, err := rc.reviewCollection.Find(ctx, bson.M{"reviewer_id": reviewId})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchQueryDB.All(ctx, &searchReviews)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchQueryDB.Close(ctx)
		if err := searchQueryDB.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchReviews)
	}
}
