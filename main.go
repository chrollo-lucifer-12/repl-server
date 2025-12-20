package main

import (
	"github.com/chrollo-lucifer-12/repl/db"
	"github.com/chrollo-lucifer-12/repl/docker"
	"github.com/chrollo-lucifer-12/repl/env"
	"github.com/chrollo-lucifer-12/repl/logger"
	"github.com/chrollo-lucifer-12/repl/server"
)

func main() {
	l := logger.NewSlogLogger()
	d := docker.NewDockerClient()
	e := env.Load()
	if e == nil {
		l.Error("no env")
	}
	db := db.NewDB(e, l)
	s := server.NewServer(l, d, db)
	err := s.Start()
	if err != nil {
		l.Error("error starting server ", err)
	}
}
