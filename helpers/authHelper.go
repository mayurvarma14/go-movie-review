package helpers

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func VerifyUserType(c *gin.Context, role string) error {
	userType := c.GetString("user_type")
	if userType != role {
		return fmt.Errorf("unauthorized: user is not a %s", role)
	}
	return nil
}

func MatchUserID(c *gin.Context, userID string) error {
	uid := c.GetString("uid")
	if uid != userID {
		return errors.New("unauthorized: user ID mismatch")
	}
	return nil
}

func MaskPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(bytes), nil
}

func ConfirmPassword(hashedPassword, userPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))
	if err != nil {
		log.Printf("Error comparing password: %v", err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, errors.New("incorrect password")
		}
		return false, fmt.Errorf("comparing password: %w", err)
	}
	return true, nil
}
