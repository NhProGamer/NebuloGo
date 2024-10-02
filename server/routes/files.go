package routes

import (
	"NebuloGo/config"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"io"
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

	// Construire le chemin absolu pour le chemin demandé
	absRequestedPath, err := filepath.Abs(cleanedPath)
	if err != nil {
		return false
	}

	// Construire le chemin absolu basé sur le répertoire de base
	absBaseDir, err := filepath.Abs(baseDir)
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
		userPath := filepath.Join(config.Configuration.Storage.Directory, claims["user_id"].(string))
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
	//filename := c.DefaultQuery("filename", "")
	claims := jwt.ExtractClaims(c)
	//c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1024<<20)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16*1024*1024*1024)

	if accessManager(c, requestedUserID) {
		userPath := filepath.Join(config.Configuration.Storage.Directory, claims["user_id"].(string))
		filePath := filepath.Join(userPath, path)
		if !isPathAllowed(userPath, filePath) {
			c.String(http.StatusForbidden, "Accès refusé")
			return
		}
		/*file, err := c.FormFile("file")
		if err != nil {
			fmt.Println(err)
			c.String(400, "Failed to get file: %s", err.Error())
			return
		}*/
		//filepath.Join(filePath, file.Filename)
		// Créer le fichier sur le disque
		//out, err := os.Create(filepath.Join(filePath, file.Filename))

		// Récupérer le fichier à partir du formulaire
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "Erreur lors de la récupération du fichier: %s", err.Error())
			return
		}
		defer file.Close()

		// Chemin complet pour enregistrer le fichier
		dst := filepath.Join(filePath, header.Filename)

		// Créer le fichier sur le disque
		out, err := os.Create(dst)
		if err != nil {
			c.String(http.StatusInternalServerError, "Erreur lors de la création du fichier: %s", err.Error())
			return
		}
		defer func() {
			if err := out.Close(); err != nil {
				fmt.Printf("Erreur lors de la fermeture du fichier: %v\n", err)
			}
		}()

		// Copier le contenu du fichier téléchargé directement sur le disque
		if _, err := io.Copy(out, file); err != nil {
			// Si l'écriture échoue, supprimer le fichier
			os.Remove(dst)
			c.String(http.StatusInternalServerError, "Erreur lors de l'écriture du fichier: %s", err.Error())
			return
		}

		c.String(http.StatusOK, "Fichier uploadé avec succès: %s", header.Filename)

	}
}

func DeleteFile(c *gin.Context) {
	requestedUserID := c.DefaultQuery("userId", "")
	path := c.DefaultQuery("path", "")
	filename := c.DefaultQuery("filename", "")
	claims := jwt.ExtractClaims(c)

	if accessManager(c, requestedUserID) {
		userPath := filepath.Join(config.Configuration.Storage.Directory, claims["user_id"].(string))
		filePath := filepath.Join(userPath, path, filename)
		if !isPathAllowed(userPath, filePath) {
			c.String(http.StatusForbidden, "Accès refusé")
			return
		}
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

	} else {
		c.JSON(http.StatusForbidden, "Forbidden")
	}
}

func MoveFile(c *gin.Context) {
	requestedUserID := c.DefaultQuery("userId", "")
	path := c.DefaultQuery("path", "")
	newpath := c.DefaultQuery("newpath", "")
	filename := c.DefaultQuery("filename", "")
	newFilename := c.DefaultQuery("newName", "")
	claims := jwt.ExtractClaims(c)

	if accessManager(c, requestedUserID) {
		userPath := filepath.Join(config.Configuration.Storage.Directory, claims["user_id"].(string))
		filePath := filepath.Join(userPath, path, filename)
		newFilePath := ""
		if newpath == "" {
			newFilePath = filepath.Join(userPath, path, newFilename)
		} else {
			newFilePath = filepath.Join(userPath, newpath, newFilename)
		}
		if !isPathAllowed(userPath, filePath) {
			c.String(http.StatusForbidden, "Accès refusé")
			return
		}
		if !isPathAllowed(userPath, newFilePath) {
			c.String(http.StatusForbidden, "Accès refusé")
			return
		}

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, "Not Found")
			return
		}
		err := os.Rename(filePath, newFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "")
		} else {
			c.JSON(http.StatusOK, "")
		}

	} else {
		c.JSON(http.StatusForbidden, "Forbidden")
	}
}

func CreateFolder(c *gin.Context) {
	requestedUserID := c.DefaultQuery("userId", "")
	path := c.DefaultQuery("path", "")
	folderName := c.DefaultQuery("folderName", "")
	claims := jwt.ExtractClaims(c)

	if accessManager(c, requestedUserID) {
		userPath := filepath.Join(config.Configuration.Storage.Directory, claims["user_id"].(string))
		folderPath := filepath.Join(userPath, path, folderName)
		if !isPathAllowed(userPath, folderPath) {
			c.String(http.StatusForbidden, "Accès refusé")
			return
		}
		err := os.Mkdir(folderPath, 0777)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "")
		} else {
			c.JSON(http.StatusOK, "")
		}
	}
}

func DeleteFolder(c *gin.Context) {
	requestedUserID := c.DefaultQuery("userId", "")
	path := c.DefaultQuery("path", "")
	folderName := c.DefaultQuery("folderName", "")
	claims := jwt.ExtractClaims(c)

	if accessManager(c, requestedUserID) {
		userPath := filepath.Join(config.Configuration.Storage.Directory, claims["user_id"].(string))
		folderPath := filepath.Join(userPath, path, folderName)
		if !isPathAllowed(userPath, folderPath) {
			c.String(http.StatusForbidden, "Accès refusé")
			return
		}
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, "Not Found")
			return
		}
		err := os.RemoveAll(folderPath) // specify the file path
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, "Ok")
		}

	} else {
		c.JSON(http.StatusForbidden, "Forbidden")
	}
}

func Content(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	requestedUserID := c.DefaultQuery("userId", "")
	path := c.DefaultQuery("path", "")
	userPath := filepath.Join(config.Configuration.Storage.Directory, claims["user_id"].(string))
	requestedPath := filepath.Join(userPath, path)

	if accessManager(c, requestedUserID) {
		if !isPathAllowed(userPath, requestedPath) {
			c.String(http.StatusForbidden, "Accès refusé")
			return
		}
		var items []interface{}
		files, err := os.ReadDir(requestedPath)
		if err != nil {
			c.JSON(http.StatusNotFound, "Not Found")
			return
		}

		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				c.String(500, err.Error())
				return
			}
			if fileInfo == nil {
				c.String(500, "No file Informations")
				return
			}
			if file.IsDir() {
				items = append(items, map[string]interface{}{"Type": "directory", "Name": file.Name(), "Time": fileInfo.ModTime()})
			} else {
				items = append(items, map[string]interface{}{"Type": "file", "Name": file.Name(), "Time": fileInfo.ModTime()})
			}
		}

		c.JSON(http.StatusOK, items)
	} else {
		c.JSON(http.StatusForbidden, "Forbidden")
	}
}
