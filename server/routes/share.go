package routes

import (
	"NebuloGo/config"
	"NebuloGo/database"
	"NebuloGo/utils"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func DownloadSharePublic(c *gin.Context) {
	requestedShareId := c.DefaultQuery("shareId", "")

	shareId, err := primitive.ObjectIDFromHex(requestedShareId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	share, err := database.ApplicationDataManager.ShareManager.GetShareFile(shareId)
	if err != nil {
		c.JSON(http.StatusNotFound, "Share not found")
		return
	}

	if !share.Public {
		c.JSON(http.StatusForbidden, "Forbidden")
		return
	}

	filePath := filepath.Join(config.Configuration.Storage.Directory, share.FilePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, "Not Found")
		return
	}

	// Serve the file
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	c.File(filePath)
}

func CreateShare(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	path := c.DefaultQuery("path", "")

	if path == "" {
		c.String(http.StatusBadRequest, "Mauvaise requête")
		return
	}

	userPath := filepath.Join(config.Configuration.Storage.Directory, claims["user_id"].(string))
	filePath := filepath.Join(userPath, path)
	if !utils.IsPathAllowed(userPath, filePath) {
		c.String(http.StatusForbidden, "Accès refusé")
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, "Not Found")
		return
	}
	userDatabaseId, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	err = database.ApplicationDataManager.ShareManager.CreateShare(userDatabaseId, filePath, []primitive.ObjectID{}, true, time.Unix(1<<63-62135596801, 999999999))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	c.JSON(http.StatusOK, "Share created!")
}
