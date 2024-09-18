package routes

import (
	"NebuloGo/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetLogout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", strings.SplitN(config.Configuration.Server.ServerURL, "//", 2)[1], strings.HasPrefix(config.Configuration.Server.ServerURL, "https://"), true)
	c.Redirect(http.StatusFound, "/login")
}
