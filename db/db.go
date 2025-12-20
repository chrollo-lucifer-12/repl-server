package db

import (
	"context"

	"github.com/chrollo-lucifer-12/repl/env"
	"github.com/chrollo-lucifer-12/repl/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
	l  logger.Logger
}

type CreatedUser struct {
	Id    uint
	Email string
}

type CreatedProject struct {
	Slug string
	Id   uint
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

func (d *DB) CreateUser(email string, password string) (*CreatedUser, error) {
	user := User{Email: email, Password: password}
	ctx := context.Background()

	result := gorm.WithResult()
	if err := gorm.G[User](d.db, result).Create(ctx, &user); err != nil {
		return nil, err
	}

	return &CreatedUser{Id: user.ID, Email: user.Email}, nil

}

func (d *DB) CreateProject(slug string, userId uint) (*CreatedProject, error) {
	project := Project{Slug: slug, UserId: userId}
	ctx := context.Background()

	result := gorm.WithResult()
	if err := gorm.G[Project](d.db, result).Create(ctx, &project); err != nil {
		return nil, err
	}

	return &CreatedProject{Slug: project.Slug, Id: project.ID}, nil
}
