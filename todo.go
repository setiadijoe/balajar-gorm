package belajargorm

import (
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Todo struct {
	ID        int                   `gorm:"primary_key;column:id;autoIncrement"`
	UserID    string                `gorm:"column:user_id"`
	Task      string                `gorm:"column:task"`
	CreatedAt int64                 `gorm:"column:created_at;autoCreateTime:nano"`
	UpdatedAt int64                 `gorm:"column:updated_at;autoUpdateTime:nano"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:nano"`
}

type TodoGorm struct {
	UserID string `gorm:"column:user_id"`
	Task   string `gorm:"column:task"`
	gorm.Model
}
