package server

import (
	"NebuloGo/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Cfg *config.Config
	Gin *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Cfg: cfg,
		Gin: gin.Default(),
	}
}

func (server *Server) Run(addr string) error {
	return server.Gin.Run(addr)
}
