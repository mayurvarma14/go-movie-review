package helpers

import (
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}
