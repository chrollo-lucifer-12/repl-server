package main

import (
	"github.com/chrollo-lucifer-12/repl/docker"
	"github.com/chrollo-lucifer-12/repl/logger"
	"github.com/chrollo-lucifer-12/repl/server"
)

func main() {
	l := logger.NewSlogLogger()
	d := docker.NewDockerClient()
	s := server.NewServer(l, d)
	err := s.Start()
	if err != nil {
		l.Error("error starting server ", err)
	}
}
