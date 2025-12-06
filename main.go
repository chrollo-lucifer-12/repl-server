package main

import (
	"github.com/chrollo-lucifer-12/repl/logger"
	"github.com/chrollo-lucifer-12/repl/server"
)

func main() {
	l := logger.NewSlogLogger()
	s := server.NewServer(l)
	err := s.Start()
	if err != nil {
		l.Error("error starting server ", err)
	}
}
