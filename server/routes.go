package server

import (
	"NebuloGo/server/routes"
)

func ConfigureRoutes(server *Server) {
	router := server.Gin
	router.LoadHTMLGlob("templates/*")
	router.GET("/", routes.GetLoginPage)

}
