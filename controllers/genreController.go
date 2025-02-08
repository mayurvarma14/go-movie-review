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

type GenreController struct {
	genreCollection *mongo.Collection
	validate        *validator.Validate
}

func NewGenreController(db *database.Database) *GenreController {
	return &GenreController{
		genreCollection: db.Client.Database(db.Name).Collection("genre"),
		validate:        validator.New(),
	}
}

func (gc *GenreController) CreateGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.VerifyUserType(c, helpers.AdminRole); err != nil {
			helpers.HandleError(c, http.StatusForbidden, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var genre models.Genre
		if err := c.BindJSON(&genre); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("binding JSON: %w", err))
			return
		}

		if err := gc.validate.Struct(&genre); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("validation: %w", err))
			return
		}

		count, err := gc.genreCollection.CountDocuments(ctx, bson.M{"name": bson.M{"$regex": *genre.Name, "$options": "i"}})
		if err != nil {
			log.Printf("Error checking genre: %v", err)
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("checking genre: %w", err))
			return
		}
		if count > 0 {
			helpers.HandleError(c, http.StatusBadRequest, errors.New("genre already exists"))
			return
		}

		genre.ID = bson.NewObjectID()
		genre.CreatedAt = time.Now()
		genre.UpdatedAt = time.Now()

		result, err := gc.genreCollection.InsertOne(ctx, genre)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("inserting genre: %w", err))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Genre created successfully", "genre_id": result.InsertedID})
	}
}

func (gc *GenreController) GetGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		genreIDStr := c.Param("genre_id")
		genreID, err := strconv.Atoi(genreIDStr)
		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid genre ID: %w", err))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var genre models.Genre
		err = gc.genreCollection.FindOne(ctx, bson.M{"genre_id": genreID}).Decode(&genre)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				helpers.HandleError(c, http.StatusNotFound, errors.New("genre not found"))
			} else {
				helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("finding genre: %w", err))
			}
			return
		}

		c.JSON(http.StatusOK, genre)
	}
}

func (gc *GenreController) GetGenres() gin.HandlerFunc {
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

		cursor, err := gc.genreCollection.Find(ctx, bson.M{}, findOptions)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("finding genres: %w", err))
			return
		}
		defer cursor.Close(ctx)

		var genres []models.Genre
		if err := cursor.All(ctx, &genres); err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("decoding genres: %w", err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"genres": genres})
	}
}

func (gc *GenreController) EditGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.VerifyUserType(c, helpers.AdminRole); err != nil {
			helpers.HandleError(c, http.StatusForbidden, err)
			return
		}

		genreIDStr := c.Param("genre_id")
		genreID, err := strconv.Atoi(genreIDStr)
		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid genre ID: %w", err))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var updatedGenre models.Genre
		if err := c.BindJSON(&updatedGenre); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("binding JSON: %w", err))
			return
		}

		if err := gc.validate.Struct(&updatedGenre); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("validation: %w", err))
			return
		}

		update := bson.M{
			"$set": bson.M{
				"name":       updatedGenre.Name,
				"updated_at": time.Now(),
			},
		}

		result, err := gc.genreCollection.UpdateOne(ctx, bson.M{"genre_id": genreID}, update)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("updating genre: %w", err))
			return
		}

		if result.MatchedCount == 0 {
			helpers.HandleError(c, http.StatusNotFound, errors.New("genre not found"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Genre updated successfully"})
	}
}

func (gc *GenreController) DeleteGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.VerifyUserType(c, helpers.AdminRole); err != nil {
			helpers.HandleError(c, http.StatusForbidden, err)
			return
		}

		genreIDStr := c.Param("genre_id")
		genreID, err := strconv.Atoi(genreIDStr)
		if err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid genre ID: %w", err))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := gc.genreCollection.DeleteOne(ctx, bson.M{"genre_id": genreID})
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("deleting genre: %w", err))
			return
		}

		if result.DeletedCount == 0 {
			helpers.HandleError(c, http.StatusNotFound, errors.New("genre not found"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Genre deleted successfully"})
	}
}
