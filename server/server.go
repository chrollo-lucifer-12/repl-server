package server

import (
	"github.com/chrollo-lucifer-12/repl/logger"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r *gin.Engine
	l logger.Logger
}

func NewServer(l logger.Logger) ServerManager {
	r := gin.Default()

	return &Server{r: r, l: l}
}

func (s *Server) Start() error {
	s.r.GET("/me", func(c *gin.Context) {
		c.Writer.Write([]byte("hi"))
	})
	s.r.GET("/ws", s.wsHandler)
	s.l.Info("server running on port :", "3000")
	err := s.r.Run(":3000")
	return err
}
