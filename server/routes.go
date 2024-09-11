package server

import (
	"NebuloGo/server/auth"
	"NebuloGo/server/routes"
)

func ConfigureRoutes(server *Server) {
	router := server.Gin
	//Serve html templates and static files
	router.LoadHTMLGlob("templates/*")
	router.Static("static", "public/static")

	//Route for login page
	router.GET("/login", routes.GetLoginPage)
	router.GET("/register", routes.GetRegisterPage)

	//Route for testing purposes
	//router.GET("/data", routes.GetDatabaseInfos)

	router.NoRoute(auth.JWTMiddleware.MiddlewareFunc(), auth.HandleNoRoute())

	//Api Routes
	router.POST("/api/v1/auth/login", auth.JWTMiddleware.LoginHandler)
	authorization := router.Group("/api/v1/auth", auth.JWTMiddleware.MiddlewareFunc())
	authorization.GET("/refresh_token", auth.JWTMiddleware.RefreshHandler)
	authorization.GET("/profile", auth.HelloHandler)
}
