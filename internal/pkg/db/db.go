package db

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TicketTransfer struct {
	gorm.Model

	TicketID   uint `gorm:"not null;index"`
	FromUserID uint `gorm:"not null;index"`
	ToUserID   uint `gorm:"not null;index"`

	Status string `gorm:"type:varchar(32);not null;default:pending;index"`

	TransferType string `gorm:"type:varchar(32);not null;default:gift;index"`

	Price float64

	Reason string `gorm:"type:text"`

	ReviewedBy uint

	ReviewedAt *time.Time
}

type MarketplaceListing struct {
	gorm.Model

	TicketID uint `gorm:"not null;index"`
	SellerID uint `gorm:"not null;index"`

	Price float64 `gorm:"not null"`

	Status string `gorm:"type:varchar(32);not null;default:active;index"`

	BuyerID uint

	Description string `gorm:"type:text"`
}
type PromoCode struct {
	gorm.Model

	Code string `gorm:"type:varchar(64);uniqueIndex;not null"`

	EventID uint `gorm:"index"`

	DiscountType string `gorm:"type:varchar(16);not null"`

	DiscountValue float64 `gorm:"not null"`

	MinAmount float64 `gorm:"default:0"`

	MaxUses int `gorm:"default:0"`

	UsedCount int `gorm:"default:0"`

	StartTime time.Time
	EndTime   time.Time

	IsActive bool `gorm:"default:true"`
}

func NewConnection(dsn string, maxOpenConns, maxIdleConns, connMaxLifetime int) (*gorm.DB, error) {
	//dsn :="host=localhost user=postgres password=123456 dbname=ticket port=5432 sslmode=disable"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	return db, nil
}

func Migrate(db *gorm.DB, models ...interface{}) error {
	if err := db.AutoMigrate(models...); err != nil {
		return err
	}

	return nil
}
