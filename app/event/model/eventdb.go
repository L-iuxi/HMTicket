package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	EventStatusDraft   = "draft"
	EventStatusReady   = "ready"
	EventStatusSelling = "selling"
	EventStatusClosed  = "closed"
)

type Event struct {
	gorm.Model

	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Location    string `gorm:"type:varchar(255);not null"`
	CoverImage  string `gorm:"type:varchar(512)"`

	StartTime  time.Time `gorm:"not null;index"`
	EndTime    time.Time `gorm:"not null;index"`
	Status     string    `gorm:"type:varchar(32);default:draft;not null;index"`
	TotalStock int       `gorm:"default:0"`
}

type TicketType struct {
	gorm.Model

	EventID    uint64  `gorm:"not null;index"`
	ShowID     uint64  `gorm:"not null;index"`
	Name       string  `gorm:"type:varchar(128);not null"`
	Price      float64 `gorm:"not null"`
	Stock      int32   `gorm:"not null;default:0"`
	MaxPerUser int32   `gorm:"not null;default:1"`
	SortOrder  int32   `gorm:"default:0"`
}
type Show struct {
	gorm.Model

	EventID uint64 `gorm:"not null;index"`

	Name string `gorm:"type:varchar(128);not null"`

	ShowTime time.Time `gorm:"not null;index"`
	EndTime  time.Time `gorm:"not null"`

	Status    string `gorm:"type:varchar(32);default:draft;not null;index"`
	Venue     string `gorm:"type:varchar(128);not null"`
	SoldCount int    `gorm:"default:0"`

	SortOrder int32 `gorm:"default:0"`
}
