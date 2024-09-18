package routes

import (
	"github.com/gin-gonic/gin"
)

func GetLoginPage(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}
