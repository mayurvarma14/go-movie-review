package helper

import (
	//"context"
	//"fmt"
	"context"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JwtSignedDetails struct {
	Email     string
	Name      string
	Username  string
	Uid       string
	User_type string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, name string, userName string, userType string, uid string) (
	signedToken string,
	signedRefreshToken string,
	err error) {
	claims := &JwtSignedDetails{
		Email:     email,
		Name:      name,
		Username:  userName,
		Uid:       uid,
		User_type: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(12)).Unix(),
		},
	}

	refreshClaims := &JwtSignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(100)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *JwtSignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtSignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*JwtSignedDetails)
	if !ok {
		msg = "This token is incorrect. Sorry!"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Ooops looks like your token has expired"
		return
	}
	return claims, msg
}

func UpdateTokens(signedToken string, signedRefreshToken string, userId string, userCollection *mongo.Collection) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateTok bson.D

	updateTok = append(updateTok, bson.E{Key: "token", Value: signedToken})
	updateTok = append(updateTok, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateTok = append(updateTok, bson.E{Key: "updated_at", Value: Updated_at})

	filter := bson.M{"user_id": userId}

	opts := options.UpdateOne().SetUpsert(true)

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateTok},
		},
		opts,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
	return
}
