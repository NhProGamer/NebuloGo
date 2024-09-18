package main

import (
	"NebuloGo/config"
	"NebuloGo/database"
	"NebuloGo/server"
	"NebuloGo/server/auth"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func main() {
	config.LoadConfig()
	sqlite.InitSqliteDB()

	if config.Configuration.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := server.NewServer()
	auth.InitJWT()
	server.ConfigureMiddlewares(app)
	server.ConfigureRoutes(app)

	err := app.Gin.SetTrustedProxies(config.Configuration.Server.TrustedProxies)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(config.Configuration.Server.Host + ":" + strconv.Itoa(config.Configuration.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
}
