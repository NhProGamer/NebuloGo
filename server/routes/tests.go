package routes

import (
	sqlite "NebuloGo/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDatabaseInfos(c *gin.Context) {
	c.JSON(http.StatusOK, sqlite.UsersCache)
}

func CreateUser(c *gin.Context) {

	c.JSON(http.StatusOK, sqlite.UsersCache)
}
