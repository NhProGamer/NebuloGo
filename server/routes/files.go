package routes

import (
	"NebuloGo/config"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"io"
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
	claims := jwt.ExtractClaims(c)
	requestedUserID := c.DefaultQuery("userId", claims["user_id"].(string))
	path := c.DefaultQuery("path", "")
	filename := c.DefaultQuery("filename", "")

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
	claims := jwt.ExtractClaims(c)
	requestedUserID := c.DefaultQuery("userId", claims["user_id"].(string))
	path := c.DefaultQuery("path", "")

	if accessManager(c, requestedUserID) {
		userPath := filepath.Join(config.Configuration.Storage.Directory, claims["user_id"].(string))
		filePath := filepath.Join(userPath, path)

		if !isPathAllowed(userPath, filePath) {
			c.String(http.StatusForbidden, "Accès refusé")
			return
		}

		// Récupérer le corps de la requête et créer un lecteur multipart
		reader, err := c.Request.MultipartReader()
		if err != nil {
			c.String(http.StatusBadRequest, "Échec de la récupération du lecteur multipart: %s", err.Error())
			return
		}

		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break // Fin des parties
			}
			if err != nil {
				c.String(http.StatusInternalServerError, "Erreur lors de la lecture des parties: %s", err.Error())
				return
			}

			// Vérifiez si la partie est un fichier
			if part.FileName() != "" {
				// Récupérer le nom du fichier
				fileName := filepath.Base(part.FileName())
				destinationPath := filepath.Join(filePath, fileName)

				// Ouvrir le fichier de destination
				outFile, err := os.Create(destinationPath)
				if err != nil {
					c.String(http.StatusInternalServerError, "Erreur lors de la création du fichier: %s", err.Error())
					return
				}
				defer outFile.Close()

				// Écrire le contenu du fichier dans le fichier de destination
				// Vous pouvez aussi directement streamer sans lire entièrement le contenu en mémoire
				if _, err := io.Copy(outFile, part); err != nil {
					c.String(http.StatusInternalServerError, "Erreur lors de l'écriture du fichier: %s", err.Error())
					return
				}
			}
		}

		c.String(http.StatusOK, "Fichiers uploadés avec succès.")
	}
}

func DeleteFile(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	requestedUserID := c.DefaultQuery("userId", claims["user_id"].(string))
	path := c.DefaultQuery("path", "")
	filename := c.DefaultQuery("filename", "")
	log.Println(filename)

	if accessManager(c, requestedUserID) {
		userPath := filepath.Join(config.Configuration.Storage.Directory, claims["user_id"].(string))
		filePath := filepath.Join(userPath, path, filename)
		if !isPathAllowed(userPath, filePath) {
			c.String(http.StatusForbidden, "Accès refusé")
			return
		}
		log.Println(filePath)
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
	claims := jwt.ExtractClaims(c)
	requestedUserID := c.DefaultQuery("userId", claims["user_id"].(string))
	path := c.DefaultQuery("path", "")
	newpath := c.DefaultQuery("newpath", "")
	filename := c.DefaultQuery("filename", "")
	newFilename := c.DefaultQuery("newName", "")

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
	claims := jwt.ExtractClaims(c)
	requestedUserID := c.DefaultQuery("userId", claims["user_id"].(string))
	path := c.DefaultQuery("path", "")
	folderName := c.DefaultQuery("folderName", "")

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
	requestedUserID := c.DefaultQuery("userId", claims["user_id"].(string))
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
				items = append(items, map[string]interface{}{"Type": "file", "Name": file.Name(), "Size": fileInfo.Size(), "Time": fileInfo.ModTime()})
			}
		}

		c.JSON(http.StatusOK, items)
	} else {
		c.JSON(http.StatusForbidden, "Forbidden")
	}
}
