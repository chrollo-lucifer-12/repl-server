package server

import (
	"github.com/chrollo-lucifer-12/repl/logger"
	"github.com/chrollo-lucifer-12/repl/terminal"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r *gin.Engine
	l logger.Logger
	t terminal.Terminal
}

func NewServer(l logger.Logger, t terminal.Terminal) ServerManager {
	r := gin.Default()

	return &Server{r: r, l: l, t: t}
}

func (s *Server) Start() error {
	s.r.GET("/me", func(c *gin.Context) {
		c.Writer.Write([]byte("hi"))
	})
	go func() {
		if err := s.t.Start(); err != nil {
			s.l.Error("terminal start error:", err)
		}
		s.l.Info("terminal started")
	}()
	s.r.GET("/ws", s.wsHandler)
	s.l.Info("server running on port :", "3000")
	err := s.r.Run(":3000")
	return err
}
