package belajargorm

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserID  string `gorm:"user_id"`
	Balance int64  `gorm:"balance"`
	User    *User  `gorm:"foreignKey:user_id;references:id"`
}
