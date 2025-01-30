package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mayurvarma14/go-movie-review/database"
	helper "github.com/mayurvarma14/go-movie-review/helpers"
	"github.com/mayurvarma14/go-movie-review/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userCollection *mongo.Collection
}

func NewUserController(db *database.Database) *UserController {
	return &UserController{
		userCollection: db.Client.Database(os.Getenv("MONGO_INITDB_DATABASE")).Collection("user"),
	}
}

var validate = validator.New()

func MaskPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func ConfirmPassword(userPassword string, passwordEntered string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordEntered), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Looks like you entered a wrong")
		check = false
	}
	return check, msg
}

func (uc *UserController) SignUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		//Check to see if data being passed meets the requirements
		if validationError := validate.Struct(&user); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()}})
			return
		}

		//Check to see if name exists
		regexMatch := bson.M{"$regex": bson.Regex{Pattern: *user.Email, Options: "i"}}
		emailCount, emailErr := uc.userCollection.CountDocuments(ctx, bson.M{"email": regexMatch})
		usernameMatch := bson.M{"$regex": bson.Regex{Pattern: *user.Username, Options: "i"}}
		usernameCount, usernameErr := uc.userCollection.CountDocuments(ctx, bson.M{"username": usernameMatch})
		defer cancel()
		if emailErr != nil {
			log.Panic(emailErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occurred while checking for this email"})
		}
		if emailCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Looks like this email already exists", "count": emailCount})
			return
		}
		if usernameErr != nil {
			log.Panic(usernameErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occurred while checking for this email / username"})
		}
		if usernameCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Looks like this username already exists", "count": usernameCount})
			return
		}

		//To hash the password before sending it to the db
		password := MaskPassword(*user.Password)
		user.Password = &password
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = bson.NewObjectID()
		user.User_id = user.ID.Hex()

		//Sign details to token
		token, refreshToken, _ := helper.GenerateAllTokens(
			*user.Email,
			*user.Name,
			*user.Username,
			*user.User_type,
			user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		//To add a new user to the database
		newUser := models.User{
			ID:            user.ID,
			User_id:       user.ID.Hex(),
			Name:          user.Name,
			Username:      user.Username,
			Email:         user.Email,
			Password:      user.Password,
			Created_at:    user.Created_at,
			Updated_at:    user.Updated_at,
			Token:         user.Token,
			User_type:     user.User_type,
			Refresh_token: user.Refresh_token,
		}

		result, err := uc.userCollection.InsertOne(ctx, newUser)

		//Error messages
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"Status":  http.StatusCreated,
			"Message": "User created successfully!",
			"Data":    map[string]interface{}{"data": result}})
	}

}

func (uc *UserController) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var retrievedUser models.User
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := uc.userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&retrievedUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "your email or password is incorrect"})
			return
		}

		passwordIsValid, msg := ConfirmPassword(*user.Password, *retrievedUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if retrievedUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "oops account not found"})
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*retrievedUser.Email, *retrievedUser.Name, *retrievedUser.Username, *retrievedUser.User_type, retrievedUser.User_id)
		helper.UpdateTokens(token, refreshToken, retrievedUser.User_id, uc.userCollection)
		err = uc.userCollection.FindOne(ctx, bson.M{"user_id": retrievedUser.User_id}).Decode(&retrievedUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, retrievedUser)
	}
}

func (uc *UserController) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := uc.userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func (uc *UserController) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.VerifyUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
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
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		match := bson.D{{Key: "$match", Value: bson.D{{}}}}
		group := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "Total number", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "Total number", Value: 1},
				{Key: "All users", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := uc.userCollection.Aggregate(ctx, mongo.Pipeline{
			match, group, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ooops something went wrong. Can't fetch all users"})
		}
		var totalUsers []bson.M
		if err = result.All(ctx, &totalUsers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, totalUsers[0])
	}
}
