package db

import (
	"github.com/chrollo-lucifer-12/repl/env"
	"github.com/chrollo-lucifer-12/repl/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
	l  logger.Logger
}

func NewDB(env *env.Env, l logger.Logger) *DB {
	db, err := gorm.Open(postgres.Open(env.DSN), &gorm.Config{})
	if err != nil {
		return nil
	}

	if err := db.AutoMigrate(&User{}, &Project{}); err != nil {
		l.Error("error migrating db", err.Error())
		return nil
	}

	d := &DB{db: db, l: l}
	return d
}
