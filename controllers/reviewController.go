package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mayurvarma14/go-movie-review/database"
	"github.com/mayurvarma14/go-movie-review/helpers"
	"github.com/mayurvarma14/go-movie-review/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReviewController struct {
	reviewCollection *mongo.Collection
	validate         *validator.Validate
}

func NewReviewController(db *database.Database) *ReviewController {
	return &ReviewController{
		reviewCollection: db.Client.Database(db.Name).Collection("review"),
		validate:         validator.New(),
	}
}

func (rc *ReviewController) AddReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType := c.GetString("user_type")
		if userType != helpers.UserRole {
			helpers.HandleError(c, http.StatusForbidden, errors.New("only users can add reviews"))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var review models.Reviews
		if err := c.BindJSON(&review); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("binding JSON: %w", err))
			return
		}

		if err := rc.validate.Struct(&review); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("validation: %w", err))
			return
		}

		reviewerID := c.GetString("uid") // Get reviewer ID from JWT
		if reviewerID == "" {
			helpers.HandleError(c, http.StatusInternalServerError, errors.New("reviewer ID not found in token"))
			return
		}
		objectReviewerID, err := bson.ObjectIDFromHex(reviewerID)
		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, errors.New("invalid reviewer ID format"))
			return
		}

		review.ID = bson.NewObjectID()
		review.ReviewerID = objectReviewerID // Use the extracted ID
		review.CreatedAt = time.Now()
		review.UpdatedAt = time.Now()

		_, err = rc.reviewCollection.InsertOne(ctx, review)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("inserting review: %w", err))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Review added successfully"})
	}
}

func (rc *ReviewController) ViewAMovieReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieIDStr := c.Query("movie_id")
		movieID, err := strconv.Atoi(movieIDStr)
		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid movie ID: %w", err))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var reviews []models.Reviews
		cursor, err := rc.reviewCollection.Find(ctx, bson.M{"movie_id": movieID})
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("finding reviews: %w", err))
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &reviews); err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("decoding reviews: %w", err))
			return
		}

		c.JSON(http.StatusOK, reviews)
	}
}

func (rc *ReviewController) DeleteReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		reviewIDStr := c.Param("id")
		reviewID, err := bson.ObjectIDFromHex(reviewIDStr) // Use ParseObjectID
		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid review ID format: %w", err))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Find the review to check ownership
		var review models.Reviews
		err = rc.reviewCollection.FindOne(ctx, bson.M{"_id": reviewID}).Decode(&review)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				helpers.HandleError(c, http.StatusNotFound, errors.New("review not found"))
			} else {
				helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("finding review: %w", err))
			}
			return
		}

		// Check if the user deleting the review is the owner or an admin.
		reviewerID := c.GetString("uid")
		objectReviewerID, err := bson.ObjectIDFromHex(reviewerID) // Use ParseObjectID
		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid reviewer ID: %w", err))
			return
		}
		if review.ReviewerID != objectReviewerID && c.GetString("user_type") != helpers.AdminRole {
			helpers.HandleError(c, http.StatusForbidden, errors.New("unauthorized to delete this review"))
			return
		}

		result, err := rc.reviewCollection.DeleteOne(ctx, bson.M{"_id": reviewID})
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("deleting review: %w", err))
			return
		}

		if result.DeletedCount == 0 {
			helpers.HandleError(c, http.StatusNotFound, errors.New("review not found")) // Should not happen, but check anyway
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
	}
}

func (rc *ReviewController) AllUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		reviewerID := c.Param("reviewer_id")
		if err := helpers.MatchUserID(c, reviewerID); err != nil {
			helpers.HandleError(c, http.StatusUnauthorized, err) // Enforce ownership
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		objectReviewerID, err := bson.ObjectIDFromHex(reviewerID) // Use ParseObjectID
		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid reviewer ID format: %w", err))
			return
		}
		var reviews []models.Reviews
		cursor, err := rc.reviewCollection.Find(ctx, bson.M{"reviewer_id": objectReviewerID})
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("finding reviews: %w", err))
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &reviews); err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("decoding reviews: %w", err))
			return
		}

		c.JSON(http.StatusOK, reviews)
	}
}
