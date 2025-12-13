package server

import (
	"github.com/chrollo-lucifer-12/repl/docker"
	"github.com/chrollo-lucifer-12/repl/logger"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r *gin.Engine
	l logger.Logger
	d *docker.DockerClient
}

func NewServer(l logger.Logger, d *docker.DockerClient) ServerManager {
	r := gin.Default()

	return &Server{r: r, l: l, d: d}
}

func (s *Server) Start() error {
	s.r.GET("/ws", s.wsHandler)
	s.l.Info("server running on port :", "3000")
	err := s.r.Run(":3000")
	return err
}
