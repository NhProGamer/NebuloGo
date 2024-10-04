package routes

import (
	"github.com/gin-gonic/gin"
)

func NebuloGoApp(c *gin.Context) {
	c.HTML(200, "app.html", gin.H{})

}
