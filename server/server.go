package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	Gin *gin.Engine
}

func NewServer() *Server {
	return &Server{
		Gin: gin.Default(),
	}
}

func (server *Server) Run(addr string) error {
	return server.Gin.Run(addr)
}
