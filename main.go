package main

import (
	"NebuloGo/config"
	"NebuloGo/database"
	"NebuloGo/server"
	"NebuloGo/server/auth"
	"log"
	"strconv"
)

func main() {
	config.LoadConfig()
	sqlite.InitSqliteDB()
	app := server.NewServer()
	auth.InitJWT()
	server.ConfigureMiddlewares(app)
	server.ConfigureRoutes(app)

	err := app.Run(config.Configuration.Host + ":" + strconv.Itoa(config.Configuration.Port))
	if err != nil {
		log.Fatal(err)
	}
}
