package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(c *gin.Context) {
	requestedUserID := c.Param("userID")
	filename := c.Param("filename")
	claims := jwt.ExtractClaims(c)

	if claims["user_id"].(string) != requestedUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	// Build file path
	filePath := filepath.Join("storage", claims["user_id"].(string), filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Serve the file
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.File(filePath)
}
