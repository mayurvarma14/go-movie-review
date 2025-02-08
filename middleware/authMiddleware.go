package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mayurvarma14/go-movie-review/helpers"
)

func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			helpers.HandleError(c, http.StatusUnauthorized, errors.New("no authorization header provided"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			helpers.HandleError(c, http.StatusUnauthorized, errors.New("invalid authorization header format. Expected 'Bearer <token>'"))
			c.Abort()
			return
		}

		clientToken := parts[1]
		if clientToken == "" {
			helpers.HandleError(c, http.StatusUnauthorized, errors.New("no authorization header provided"))
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != nil {
			helpers.HandleError(c, http.StatusUnauthorized, fmt.Errorf("validating token: %w", err))
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("username", claims.Username)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}
