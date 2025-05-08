package models

import (
	"order-service/constants"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uint                  `gorm:"primaryKey;autoIncrement"`
	UUID      uuid.UUID             `gorm:"type:uuid;not null"`
	Code      string                `gorm:"type:varchar(30);not null"`
	UserID    uuid.UUID             `gorm:"type:uuid;not null"`
	PaymentID uuid.UUID             `gorm:"type:uuid;not null"`
	Amount    float64               `gorm:"type:decimal(10,2);not null"`
	Status    constants.OrderStatus `gorm:"not null"`
	Date      time.Time             `gorm:"type:timestamp;not null"`
	IsPaid    bool                  `gorm:"type:boolean;not null"`
	PaidAt    *time.Time            `gorm:"type:timestamp"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
