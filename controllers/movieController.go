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

type MovieController struct {
	movieCollection *mongo.Collection
}

func NewMovieController(db *database.Database) *MovieController {
	return &MovieController{
		movieCollection: db.Client.Database(os.Getenv("MONGO_INITDB_DATABASE")).Collection("movie"),
	}
}

func (mc *MovieController) CreateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.VerifyUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var movie models.Movie
		defer cancel()

		if err := c.BindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		regexMatch := bson.M{"$regex": bson.Regex{Pattern: *movie.Name, Options: "i"}}
		count, err := mc.movieCollection.CountDocuments(ctx, bson.M{"name": regexMatch})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occurred while checking for the movie name"})
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this movie name already exists", "count": count})
			return
		}

		if validationError := validate.Struct(&movie); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()}})
			return
		}

		newMovie := models.Movie{
			Id:       bson.NewObjectID(),
			Name:     movie.Name,
			Topic:    movie.Topic,
			Genre_id: movie.Genre_id,

			Movie_URL:  movie.Movie_URL,
			Created_at: movie.Created_at,
			Updated_at: movie.Updated_at,
			Movie_id:   movie.Movie_id,
		}
		result, err := mc.movieCollection.InsertOne(ctx, newMovie)

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

func (mc *MovieController) GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		movieId := c.Param("movie_id")
		var movie models.Movie
		defer cancel()

		objId, _ := bson.ObjectIDFromHex(movieId)

		err := mc.movieCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    map[string]interface{}{"data": movie}})
	}
}

func (mc *MovieController) GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
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
				{Key: "movie_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := mc.movieCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching movies "})
		}
		var allMovies []bson.M
		if err = result.All(ctx, &allMovies); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMovies[0])
	}
}

func (mc *MovieController) UpdateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		movieId := c.Param("movie_id")
		var movie models.Movie
		defer cancel()
		objId, _ := bson.ObjectIDFromHex(movieId)

		if err := c.BindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationError := validate.Struct(&movie); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()}})
			return
		}

		update := bson.M{
			"name":      movie.Name,
			"topic":     movie.Topic,
			"genre_id":  movie.Genre_id,
			"movie_url": movie.Movie_URL}
		filterByID := bson.M{"_id": bson.M{"$eq": objId}}
		result, err := mc.movieCollection.UpdateOne(ctx, filterByID, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedMovie models.Movie
		if result.MatchedCount == 1 {
			err := mc.movieCollection.FindOne(ctx, filterByID).Decode(&updatedMovie)
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
			"Message": "movie updated successfully!",
			"Data":    updatedMovie})
	}
}

func (mc *MovieController) SearchMovieByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchMovies []models.Movie
		queryParam := c.Query("name")
		if queryParam == "" {
			log.Println("name is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search parameter"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchQueryDB, err := mc.movieCollection.Find(ctx, bson.M{"name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(404, "something went wrong")
			return
		}
		err = searchQueryDB.All(ctx, &searchMovies)
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
		c.IndentedJSON(200, searchMovies)
	}
}

func (mc *MovieController) SearchMovieByGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchByGenre []models.Movie
		genreId := c.Query("genre_id")
		if genreId == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		i, err := strconv.Atoi(genreId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
			log.Panic(err)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchDB, err := mc.movieCollection.Find(ctx, bson.M{"genre_id": i})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchDB.All(ctx, &searchByGenre)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchDB.Close(ctx)
		if err := searchDB.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchByGenre)
	}
}

func (mc *MovieController) DeleteMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		movieId := c.Param("movie_id")
		defer cancel()
		i, err := strconv.Atoi(movieId)
		if err != nil {
			// Handle error
		}
		result, err := mc.movieCollection.DeleteOne(ctx, bson.M{"movie_id": i})
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
					" Data":    map[string]interface{}{"data": "Movie with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			gin.H{
				"Status":  http.StatusOK,
				"Message": "success",
				"Data":    map[string]interface{}{"data": "Movie successfully deleted!"}},
		)
	}
}
