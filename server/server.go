package server

import (
	"github.com/chrollo-lucifer-12/repl/db"
	"github.com/chrollo-lucifer-12/repl/docker"
	"github.com/chrollo-lucifer-12/repl/logger"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r  *gin.Engine
	l  logger.Logger
	d  *docker.DockerClient
	db *db.DB
}

func NewServer(l logger.Logger, d *docker.DockerClient, db *db.DB) ServerManager {
	r := gin.Default()

	return &Server{r: r, l: l, d: d, db: db}
}

func (s *Server) Start() error {
	s.r.POST("/create-project", s.CreateProjectHandler)
	s.r.POST("/register", s.RegisterHandler)
	s.r.GET("/ws", s.wsHandler)
	s.l.Info("server running on port :", "3000")
	err := s.r.Run(":3000")
	return err
}
