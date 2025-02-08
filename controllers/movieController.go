package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MovieController struct {
	movieCollection *mongo.Collection
	validate        *validator.Validate
}

func NewMovieController(db *database.Database) *MovieController {
	return &MovieController{
		movieCollection: db.Client.Database(db.Name).Collection("movie"),
		validate:        validator.New(),
	}
}

func (mc *MovieController) CreateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.VerifyUserType(c, helpers.AdminRole); err != nil {
			helpers.HandleError(c, http.StatusForbidden, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.BindJSON(&movie); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("binding JSON: %w", err))
			return
		}

		if err := mc.validate.Struct(&movie); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("validation: %w", err))
			return
		}

		count, err := mc.movieCollection.CountDocuments(ctx, bson.M{"name": bson.M{"$regex": *movie.Name, "$options": "i"}})
		if err != nil {
			log.Printf("Error checking movie: %v", err)
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("checking movie: %w", err))
			return
		}
		if count > 0 {
			helpers.HandleError(c, http.StatusBadRequest, errors.New("movie already exists"))
			return
		}

		movie.ID = bson.NewObjectID()
		movie.CreatedAt = time.Now()
		movie.UpdatedAt = time.Now()

		result, err := mc.movieCollection.InsertOne(ctx, movie)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("inserting movie: %w", err))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Movie created successfully", "movie_id": result.InsertedID})
	}
}

func (mc *MovieController) GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieIDStr := c.Param("movie_id")
		movieID, err := strconv.Atoi(movieIDStr)
		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid movie ID: %w", err))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var movie models.Movie
		err = mc.movieCollection.FindOne(ctx, bson.M{"movie_id": movieID}).Decode(&movie)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				helpers.HandleError(c, http.StatusNotFound, errors.New("movie not found"))
			} else {
				helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("finding movie: %w", err))
			}
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func (mc *MovieController) GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			helpers.HandleError(c, http.StatusBadRequest, errors.New("invalid page number"))
			return
		}
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit < 1 {
			helpers.HandleError(c, http.StatusBadRequest, errors.New("invalid limit number"))
			return
		}
		skip := (page - 1) * limit

		findOptions := options.Find()
		findOptions.SetSkip(int64(skip))
		findOptions.SetLimit(int64(limit))

		cursor, err := mc.movieCollection.Find(ctx, bson.M{}, findOptions)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("finding movies: %w", err))
			return
		}
		defer cursor.Close(ctx)

		var movies []models.Movie
		if err := cursor.All(ctx, &movies); err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("decoding movies: %w", err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"movies": movies})
	}
}

func (mc *MovieController) UpdateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.VerifyUserType(c, helpers.AdminRole); err != nil {
			helpers.HandleError(c, http.StatusForbidden, err)
			return
		}

		movieIDStr := c.Param("movie_id")
		movieID, err := strconv.Atoi(movieIDStr)
		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid movie ID: %w", err))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var updatedMovie models.Movie
		if err := c.BindJSON(&updatedMovie); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("binding JSON: %w", err))
			return
		}
		if err := mc.validate.Struct(&updatedMovie); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("validation error: %w", err))
			return
		}

		update := bson.M{
			"$set": bson.M{
				"name":       updatedMovie.Name,
				"topic":      updatedMovie.Topic,
				"genre_id":   updatedMovie.GenreID,
				"movie_url":  updatedMovie.MovieURL,
				"updated_at": time.Now(),
			},
		}

		result, err := mc.movieCollection.UpdateOne(ctx, bson.M{"movie_id": movieID}, update)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("updating movie: %w", err))
			return
		}

		if result.MatchedCount == 0 {
			helpers.HandleError(c, http.StatusNotFound, errors.New("movie not found"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully"})
	}
}

func (mc *MovieController) SearchMovieByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("name")
		if query == "" {
			helpers.HandleError(c, http.StatusBadRequest, errors.New("search query parameter 'name' is required"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var movies []models.Movie
		cursor, err := mc.movieCollection.Find(ctx, bson.M{"name": bson.M{"$regex": query, "$options": "i"}})
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("searching movies: %w", err))
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &movies); err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("decoding movies: %w", err))
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

func (mc *MovieController) SearchMovieByGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		genreIDStr := c.Query("genre_id")
		genreID, err := strconv.Atoi(genreIDStr)
		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid genre ID: %w", err))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var movies []models.Movie
		cursor, err := mc.movieCollection.Find(ctx, bson.M{"genre_id": genreID})
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("searching movies by genre: %w", err))
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &movies); err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("decoding movies: %w", err))
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

func (mc *MovieController) DeleteMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.VerifyUserType(c, helpers.AdminRole); err != nil {
			helpers.HandleError(c, http.StatusForbidden, err)
			return
		}

		movieIDStr := c.Param("movie_id")
		movieID, err := strconv.Atoi(movieIDStr)

		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid movie ID: %w", err))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := mc.movieCollection.DeleteOne(ctx, bson.M{"movie_id": movieID})
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("deleting movie: %w", err))
			return
		}

		if result.DeletedCount == 0 {
			helpers.HandleError(c, http.StatusNotFound, errors.New("movie not found"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
	}
}
