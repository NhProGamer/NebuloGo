package server

import (
	"NebuloGo/server/auth"
	"NebuloGo/server/routes"
)

func ConfigureRoutes(server *Server) {
	router := server.Gin
	router.LoadHTMLGlob("templates/*")
	router.Static("static", "public/static")
	router.GET("/login", routes.GetLoginPage)
	router.GET("/data", routes.GetDatabaseInfos)

	router.POST("/api/v1/login", auth.JWTMiddleware.LoginHandler)
	router.NoRoute(auth.JWTMiddleware.MiddlewareFunc(), auth.HandleNoRoute())
	authorization := router.Group("/auth", auth.JWTMiddleware.MiddlewareFunc())
	authorization.GET("/refresh_token", auth.JWTMiddleware.RefreshHandler)
	authorization.GET("/hello", auth.HelloHandler)
}
