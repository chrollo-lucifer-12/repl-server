package main

import (
	"github.com/chrollo-lucifer-12/repl/logger"
	"github.com/chrollo-lucifer-12/repl/server"
	"github.com/chrollo-lucifer-12/repl/terminal"
)

func main() {
	l := logger.NewSlogLogger()
	t := terminal.NewBashTerminal()
	s := server.NewServer(l, t)
	err := s.Start()
	if err != nil {
		l.Error("error starting server ", err)
	}
}
