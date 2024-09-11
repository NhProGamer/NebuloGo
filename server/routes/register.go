package routes

import "github.com/gin-gonic/gin"

func GetRegisterPage(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}
