package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(c *gin.Context) {
	requestedUserID := c.DefaultQuery("userId", "")
	path := c.DefaultQuery("path", "")
	filename := c.DefaultQuery("filename", "")
	claims := jwt.ExtractClaims(c)

	if claims["user_id"].(string) != requestedUserID {
		c.JSON(http.StatusForbidden, "Forbidden")
		return
	}
	// Build file path
	filePath := filepath.Join("storage", claims["user_id"].(string), path, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, "Not Found")
		return
	}

	// Serve the file
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.File(filePath)
}

func UploadFile(c *gin.Context) {
	requestedUserID := c.DefaultQuery("UserId", "")
	path := c.DefaultQuery("path", "")
	//filename := c.Param("filename")
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
