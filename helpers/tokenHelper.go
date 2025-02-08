package helpers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JwtSignedDetails struct {
	Email    string
	Name     string
	Username string
	Uid      string
	UserType string
	jwt.StandardClaims
}

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func GenerateAllTokens(email, name, userName, userType, uid string) (string, string, error) {
	claims := &JwtSignedDetails{
		Email:    email,
		Name:     name,
		Username: userName,
		Uid:      uid,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(), // 15 minutes for access token
		},
	}

	refreshClaims := &JwtSignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 24 hours for refresh token
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		return "", "", fmt.Errorf("generating token: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secretKey)
	if err != nil {
		return "", "", fmt.Errorf("generating refresh token: %w", err)
	}

	return token, refreshToken, nil
}

func ValidateToken(signedToken string) (*JwtSignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtSignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	if claims, ok := token.Claims.(*JwtSignedDetails); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func UpdateTokens(signedToken, signedRefreshToken, userId string, userCollection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"token":         signedToken,
		"refresh_token": signedRefreshToken,
		"updated_at":    time.Now(),
	}

	filter := bson.M{"user_id": userId}
	opts := options.UpdateOne().SetUpsert(true)

	_, err := userCollection.UpdateOne(ctx, filter, bson.M{"$set": update}, opts)
	if err != nil {
		return fmt.Errorf("updating tokens: %w", err)
	}

	return nil
}
