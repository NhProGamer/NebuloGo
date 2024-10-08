package server

import (
	"NebuloGo/server/auth"
	"NebuloGo/server/routes"
)

func ConfigureRoutes(server *Server) {
	router := server.Gin
	router.MaxMultipartMemory = 8192 << 20
	//Serve html templates and static files
	router.LoadHTMLGlob("templates/*")
	router.Static("static", "public/static")

	//Routes for login system
	router.GET("/login", routes.GetLoginPage)
	router.GET("/register", routes.GetRegisterPage)
	router.GET("/logout", routes.GetLogout)

	//Routes for app
	app := router.Group("/drive", auth.JWTMiddleware.MiddlewareFunc())
	app.GET("/", routes.Drive)
	app.GET("/shares", routes.Shares)

	router.NoRoute(auth.JWTMiddleware.MiddlewareFunc(), auth.HandleNoRoute())

	//Api Routes
	router.POST("/api/v1/auth/login", auth.JWTMiddleware.LoginHandler)
	authorization := router.Group("/api/v1/auth", auth.JWTMiddleware.MiddlewareFunc())
	authorization.GET("/refresh_token", auth.JWTMiddleware.RefreshHandler)

	filesApi := router.Group("/api/v1/files", auth.JWTMiddleware.MiddlewareFunc())
	filesApi.GET("/content", routes.Content)
	filesApi.GET("/", routes.DownloadFile)
	filesApi.POST("/", routes.UploadFile)
	filesApi.PATCH("/", routes.MoveFile)
	filesApi.DELETE("/", routes.DeleteFile)
	filesApi.POST("/folder", routes.CreateFolder)
	filesApi.DELETE("/folder", routes.DeleteFolder)

	sharesApi := router.Group("/api/v1/share")
	sharesApi.POST("/", auth.JWTMiddleware.MiddlewareFunc(), routes.CreateShare)
	sharesApi.GET("/download", routes.DownloadSharePublic)
	sharesApi.GET("/", auth.JWTMiddleware.MiddlewareFunc(), routes.ListShares)
	sharesApi.DELETE("/", auth.JWTMiddleware.MiddlewareFunc(), routes.DeleteShare)
}
