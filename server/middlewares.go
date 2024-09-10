package server

import (
	"NebuloGo/server/auth"
)

func ConfigureMiddlewares(server *Server) {
	server.Gin.Use(auth.HandlerMiddleWare(auth.JWTMiddleware))
}
