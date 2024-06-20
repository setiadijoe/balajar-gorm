package belajargorm

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserID  string `gorm:"user_id"`
	Balance int64  `gorm:"balance"`
}
