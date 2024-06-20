package belajargorm

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID        int64          `gorm:"primary_key;column:id;autoIncrement"`
	UserID    string         `gorm:"column:user_id"`
	Address   string         `gorm:"column:address"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	User      User           `gorm:"foreignKey:user_id;references:id"`
}
