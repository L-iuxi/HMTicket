package model

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model

	UserID       uint `gorm:"not null;index"`
	EventID      uint `gorm:"not null;index"`
	ShowID       uint `gorm:"index"`
	TicketTypeID uint `gorm:"not null;index"`

	OrderNo string `gorm:"type:varchar(64);uniqueIndex;not null"`

	Quantity   int `gorm:"not null;default:1"`
	TotalPrice float64

	Status string `gorm:"type:varchar(32);not null;default:reserved;index"`

	QRCode string `gorm:"type:text"`

	DiscountCode string `gorm:"type:varchar(64);index"`

	RealName string `gorm:"type:varchar(64);index"`
	IDCard   string `gorm:"type:varchar(32);index"`
	Phone    string `gorm:"type:varchar(20);index"`

	TransferredTo uint `gorm:"index"`

	TransferStatus string `gorm:"type:varchar(32);default:none"`
}
