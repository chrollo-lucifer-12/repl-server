package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
	Projects []Project
}

type Project struct {
	gorm.Model
	Slug   string `gorm:"unique"`
	UserId uint
}
