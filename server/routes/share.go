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

	filePath := filepath.Join(claims["user_id"].(string), path)
	if !utils.IsPathAllowed(claims["user_id"].(string), filePath) {
		c.String(http.StatusForbidden, "Accès refusé")
		return
	}

	if _, err := os.Stat(filepath.Join(config.Configuration.Storage.Directory, filePath)); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, "Not Found")
		return
	}
	userDatabaseId, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	shareId, err := database.ApplicationDataManager.ShareManager.CreateShare(userDatabaseId, filePath, []primitive.ObjectID{}, true, time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	c.JSON(http.StatusOK, shareId.Hex())
}
