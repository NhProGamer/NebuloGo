package routes

import "github.com/gin-gonic/gin"

func Shares(c *gin.Context) {
	c.HTML(200, "shares.html", gin.H{})
}
