package model

import "gorm.io/gorm"

const (
	OrderUnpaid   = "unpaid"
	OrderPaid     = "paid"
	OrderCanceled = "cancelled"
	OrderRefunded = "refunded"
)

type Order struct {
	gorm.Model

	UserID uint `gorm:"not null;index"`

	EventID      uint64 `gorm:"not null;index"`
	ShowID       uint   `gorm:"index"`
	TicketTypeID uint   `gorm:"not null;index"`
	OrderNo      string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Quantity     int    `gorm:"not null;default:1"`
	TotalPrice   float64
	RequsetId    string `gorm:"not null;default:1"`
	Status       string `gorm:"not null;default:1"`
}
