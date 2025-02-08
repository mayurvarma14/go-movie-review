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

type UserController struct {
	userCollection *mongo.Collection
	validate       *validator.Validate
}

func NewUserController(db *database.Database) *UserController {
	return &UserController{
		userCollection: db.Client.Database(db.Name).Collection("user"),
		validate:       validator.New(),
	}
}

func (uc *UserController) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("binding JSON: %w", err))
			return
		}

		if err := uc.validate.Struct(&user); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("validation: %w", err))
			return
		}

		emailCount, err := uc.userCollection.CountDocuments(ctx, bson.M{"email": bson.M{"$regex": *user.Email, "$options": "i"}})
		if err != nil {
			log.Printf("Error checking email: %v", err)
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("checking email: %w", err))
			return
		}
		if emailCount > 0 {
			helpers.HandleError(c, http.StatusBadRequest, errors.New("email already exists"))
			return
		}

		usernameCount, err := uc.userCollection.CountDocuments(ctx, bson.M{"username": bson.M{"$regex": *user.Username, "$options": "i"}})
		if err != nil {
			log.Printf("Error checking username: %v", err)
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("checking username: %w", err))
			return
		}
		if usernameCount > 0 {
			helpers.HandleError(c, http.StatusBadRequest, errors.New("username already exists"))
			return
		}

		hashedPassword, err := helpers.MaskPassword(*user.Password)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, err)
			return
		}
		user.Password = &hashedPassword
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.ID = bson.NewObjectID()
		user.UserID = user.ID.Hex()

		token, refreshToken, err := helpers.GenerateAllTokens(*user.Email, *user.Name, *user.Username, *user.UserType, user.UserID)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, err)
			return
		}
		user.Token = &token
		user.RefreshToken = &refreshToken

		result, err := uc.userCollection.InsertOne(ctx, user)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("inserting user: %w", err))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": result.InsertedID})
	}
}

func (uc *UserController) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var loginUser, foundUser models.User

		if err := c.BindJSON(&loginUser); err != nil {
			helpers.HandleError(c, http.StatusBadRequest, fmt.Errorf("binding JSON: %w", err))
			return
		}

		err := uc.userCollection.FindOne(ctx, bson.M{"email": bson.M{"$regex": *loginUser.Email, "$options": "i"}}).Decode(&foundUser)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				helpers.HandleError(c, http.StatusUnauthorized, errors.New("invalid email or password"))
			} else {
				helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("finding user: %w", err))
			}
			return
		}

		passwordMatch, err := helpers.ConfirmPassword(*foundUser.Password, *loginUser.Password)
		if err != nil || !passwordMatch {
			helpers.HandleError(c, http.StatusUnauthorized, errors.New("invalid email or password"))
			return
		}

		token, refreshToken, err := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.Name, *foundUser.Username, *foundUser.UserType, foundUser.UserID)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, err)
			return
		}

		if err := helpers.UpdateTokens(token, refreshToken, foundUser.UserID, uc.userCollection); err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token, "refresh_token": refreshToken})
	}
}

func (uc *UserController) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helpers.MatchUserID(c, userId); err != nil && c.GetString("user_type") != helpers.AdminRole {
			helpers.HandleError(c, http.StatusUnauthorized, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		err := uc.userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				helpers.HandleError(c, http.StatusNotFound, errors.New("user not found"))
			} else {
				helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("finding user: %w", err))
			}
			return
		}

		user.Password = nil
		user.Token = nil
		user.RefreshToken = nil

		c.JSON(http.StatusOK, user)
	}
}

func (uc *UserController) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.VerifyUserType(c, helpers.AdminRole); err != nil {
			helpers.HandleError(c, http.StatusForbidden, err)
			return
		}

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

		cursor, err := uc.userCollection.Find(ctx, bson.M{}, findOptions)
		if err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("finding users: %w", err))
			return
		}
		defer cursor.Close(ctx)

		var users []models.User
		if err := cursor.All(ctx, &users); err != nil {
			helpers.HandleError(c, http.StatusInternalServerError, fmt.Errorf("decoding users: %w", err))
			return
		}

		for i := range users {
			users[i].Password = nil
			users[i].Token = nil
			users[i].RefreshToken = nil
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}
