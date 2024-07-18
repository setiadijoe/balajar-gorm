package belajargorm

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string    `gorm:"primary_key;column:id"`
	Password     string    `gorm:"column:password"`
	Name         Name      `gorm:"embedded"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Information  string    `gorm:"-"`
	Wallet       Wallet    `gorm:"foreignKey:user_id;references:id"`
	Addresses    []Address `gorm:"foreignKey:user_id;references:id"`
	LikeProducts []Product `gorm:"many2many:user_like_product;foreignKey:id;joinForeignKey:user_id;joinReferences:product_id"`
}

type Name struct {
	FirstName  string `gorm:"column:first_name"`
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name"`
}

func (u *User) TableName() string {
	return "users"
}

type UserLog struct {
	ID        int    `gorm:"primary_key;column:id;autoIncrement"`
	UserID    string `gorm:"column:user_id"`
	Action    string `gorm:"column:action"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:mili"`
	UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime:mili"`
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.ID == "" {
		u.ID = fmt.Sprintf("user-%d", time.Now().UTC().UnixMilli())
	}
	return nil
}
