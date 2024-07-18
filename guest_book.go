package belajargorm

import "gorm.io/gorm"

type GuestBook struct {
	Name    string `gorm:"column:name"`
	Email   string `gorm:"column:email"`
	Message string `gorm:"column:message"`
	gorm.Model
}
