package model

import "gorm.io/gorm"

type AdminModel struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique"`
	Password string
}