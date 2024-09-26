package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func accessManager(c *gin.Context, requestedUserID string) bool {
	claims := jwt.ExtractClaims(c)
	return claims["user_id"].(string) == requestedUserID
}

func isPathAllowed(baseDir, requestedPath string) bool {
	// Nettoyer le chemin demandé
	cleanedPath := filepath.Clean(requestedPath)

	// Construire le chemin absolu basé sur le répertoire de base
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return false
	}

	// Construire le chemin absolu pour le chemin demandé
	absRequestedPath, err := filepath.Abs(filepath.Join(baseDir, cleanedPath))
	if err != nil {
		return false
	}

	// Vérifier que le chemin demandé commence bien par le répertoire de base
	return strings.HasPrefix(absRequestedPath, absBaseDir)
}

func DownloadFile(c *gin.Context) {
	requestedUserID := c.DefaultQuery("userId", "")
	path := c.DefaultQuery("path", "")
	filename := c.DefaultQuery("filename", "")
	claims := jwt.ExtractClaims(c)

	if accessManager(c, requestedUserID) {
		userPath := filepath.Join("storage", claims["user_id"].(string))
		filePath := filepath.Join(userPath, path, filename)
		if !isPathAllowed(userPath, filePath) {
			c.String(http.StatusForbidden, "Accès refusé")
			return
		}

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, "Not Found")
			return
		}

		// Serve the file
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.File(filePath)

	} else {
		c.JSON(http.StatusForbidden, "Forbidden")
	}
}

func UploadFile(c *gin.Context) {
	requestedUserID := c.DefaultQuery("UserId", "")
	path := c.DefaultQuery("path", "")
	claims := jwt.ExtractClaims(c)

	if claims["user_id"].(string) != requestedUserID {
		c.JSON(http.StatusForbidden, "Forbidden")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.String(400, "Failed to get file: %s", err.Error())
		return
	}
	savePath := filepath.Join("storage", claims["user_id"].(string), path, file.Filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.String(500, "Failed to save file: %s", err.Error())
		return
	}
}

func DeleteFile(c *gin.Context) {
	requestedUserID := c.DefaultQuery("userId", "")
	path := c.DefaultQuery("path", "")
	actualName := c.DefaultQuery("actualName", "")
	claims := jwt.ExtractClaims(c)

	if claims["user_id"].(string) != requestedUserID {
		c.JSON(http.StatusForbidden, "Forbidden")
		return
	}
	// Build file path
	filePath := filepath.Join("storage", claims["user_id"].(string), path, actualName)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, "Not Found")
		return
	}

	err := os.Remove(filePath) // specify the file path
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	} else {
		c.JSON(http.StatusOK, "")
	}

}

func RenameFile(c *gin.Context) {
	requestedUserID := c.DefaultQuery("userId", "")
	path := c.DefaultQuery("path", "")
	newName := c.DefaultQuery("newName", "")
	actualName := c.DefaultQuery("actualName", "")
	claims := jwt.ExtractClaims(c)

	if claims["user_id"].(string) != requestedUserID {
		c.JSON(http.StatusForbidden, "Forbidden")
		return
	}
	// Build file path
	actualFilePath := filepath.Join("storage", claims["user_id"].(string), path, actualName)
	newFilePath := filepath.Join("storage", claims["user_id"].(string), path, newName)

	// Check if file exists
	if _, err := os.Stat(actualFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, "Not Found")
		return
	}

	err := os.Rename(actualFilePath, newFilePath) // specify the file path
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	} else {
		c.JSON(http.StatusOK, "")
	}

}

func Content(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	requestedUserID := c.DefaultQuery("userId", "")
	path := c.DefaultQuery("path", "")

	if claims["user_id"].(string) == requestedUserID {
		var items []interface{}
		files, err := os.ReadDir("./storage/" + requestedUserID + "/" + path)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, "Not Found")
			return
		}

		for _, file := range files {
			if file.IsDir() {
				items = append(items, map[string]interface{}{"Type": "directory", "Name": file.Name()})
			} else {
				items = append(items, map[string]interface{}{"Type": "file", "Name": file.Name()})
			}
		}

		c.JSON(http.StatusOK, items)
	} else {
		c.JSON(http.StatusForbidden, "Forbidden")
	}
}
