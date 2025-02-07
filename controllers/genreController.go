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
	"github.com/mayurvarma14/go-movie-review/helpers"
	helper "github.com/mayurvarma14/go-movie-review/helpers"
	"github.com/mayurvarma14/go-movie-review/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GenreController struct {
	genreCollection *mongo.Collection
}

func NewGenreController(db *database.Database) *GenreController {
	return &GenreController{
		genreCollection: db.Client.Database(os.Getenv("MONGO_INITDB_DATABASE")).Collection("genre"),
	}
}

func (gc *GenreController) CreateGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.VerifyUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var genre models.Genre
		if err := c.BindJSON(&genre); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		if validationError := validate.Struct(&genre); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()}})
			return
		}

		//Check to see if name exists
		regexMatch := bson.M{"$regex": bson.Regex{Pattern: *genre.Name, Options: "i"}}
		count, err := gc.genreCollection.CountDocuments(ctx, bson.M{"Name": regexMatch})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while checking for the genre name"})
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this genre name already exists", "count": count})
			return
		}

		newGenre := models.Genre{
			Id:         bson.NewObjectID(),
			Name:       genre.Name,
			Created_at: time.Now(),
			Updated_at: time.Now(),
			Genre_id:   genre.Genre_id,
		}
		result, err := gc.genreCollection.InsertOne(ctx, newGenre)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"Status":  http.StatusCreated,
			"Message": "genre created successfully",
			"Data":    map[string]interface{}{"data": result}})

	}
}

func (gc *GenreController) GetGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		genre_id := c.Param("genre_id")
		var genre models.Genre
		defer cancel()

		genreId, err := strconv.Atoi(genre_id)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while checking for the genre name"})
		}

		err = gc.genreCollection.FindOne(ctx, bson.M{"genre_id": genreId}).Decode(&genre)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"error": err.Error()},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    map[string]interface{}{"data": genre}})
	}
}

func (gc *GenreController) GetGenres() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "genre_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := gc.genreCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching genres "})
		}
		var genres []bson.M
		if err = result.All(ctx, &genres); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, genres[0])
	}
}

func (gc *GenreController) EditGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.VerifyUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		genreId := c.Param("genre_id")
		var genre models.Genre
		objId, _ := bson.ObjectIDFromHex(genreId)

		if err := c.BindJSON(&genre); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&genre); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": genre.Name}
		filterByID := bson.M{"_id": bson.M{"$eq": objId}}
		result, err := gc.genreCollection.UpdateOne(ctx, filterByID, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedGenre models.Genre
		if result.MatchedCount == 1 {
			err := gc.genreCollection.FindOne(ctx, filterByID).Decode(&updatedGenre)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"Status":  http.StatusInternalServerError,
					"Message": "error",
					"Data":    map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    updatedGenre})
	}
}

func (gc *GenreController) DeleteGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		genreId := c.Param("genre_id")
		i, err := strconv.Atoi(genreId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
			log.Panic(err)
			return
		}
		result, err := gc.genreCollection.DeleteOne(ctx, bson.M{"genre_id": i})
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
					" Data":    map[string]interface{}{"data": "Genre with specified ID not found!"}},
			)
			return

		}

		c.JSON(http.StatusOK,
			gin.H{
				"Status":  http.StatusOK,
				"Message": "success",
				"Data":    map[string]interface{}{"data": "Genre successfully deleted!"}},
		)
	}
}
