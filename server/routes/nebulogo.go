package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func NebuloGoApp(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	path := c.DefaultQuery("path", "")

	c.HTML(200, "app.html", gin.H{
		/*"Items":  items,*/
		"UserId":     claims["user_id"].(string),
		"ActualPath": path,
	})

}
