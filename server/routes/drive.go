package routes

import (
	"github.com/gin-gonic/gin"
)

func Drive(c *gin.Context) {
	c.HTML(200, "app.html", gin.H{})
}
