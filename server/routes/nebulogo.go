package routes

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"os"
)

func NebuloGoApp(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	var items []interface{}

	files, err := os.ReadDir("./storage/" + claims["user_id"].(string))
	if err != nil {
		fmt.Println("Erreur lors de la lecture du répertoire:", err)
	}
	for _, file := range files {
		if file.IsDir() {
			items = append(items, map[string]interface{}{"Type": "directory", "Name": file.Name()})
		} else {
			items = append(items, map[string]interface{}{"Type": "file", "Name": file.Name()})
		}
	}
	c.HTML(200, "app.html", gin.H{
		"Items": items,
	})

}
