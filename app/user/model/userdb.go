package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID   string `gorm:"type:varchar(6);uniqueIndex"`
	Username string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(128);uniqueIndex;not null"`
	Role     string `gorm:"type:varchar(32);default:user"`
	Phone    string `gorm:"type:varchar(20);uniqueIndex"`
	Gender   uint8  `gorm:"default:0"`
}
