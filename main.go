package main

import (
	"NebuloGo/config"
	"NebuloGo/database"
	"NebuloGo/server"
	"log"
	"strconv"
)

func main() {
	cfg := config.LoadConfig()
	sqlite.InitSqliteDB()
	app := server.NewServer(cfg)
	server.ConfigureRoutes(app)

	err := app.Run(cfg.HOST + ":" + strconv.Itoa(cfg.PORT))
	if err != nil {
		log.Fatal(err)
	}
}
