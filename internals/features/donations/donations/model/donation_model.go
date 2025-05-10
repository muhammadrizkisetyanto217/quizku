package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Donation struct {
	ID             uint       `gorm:"primaryKey"`
	UserID         uuid.UUID  `gorm:"type:uuid"`                           // Relasi ke users
	Amount         int        `gorm:"not null"`                            // Jumlah donasi
	Message        string     `gorm:"type:text"`                           // Pesan dari donor
	Status         int        `gorm:"default:0"`                           // 0 = pending, 1 = paid, 2 = expired, 3 = canceled
	OrderID        string     `gorm:"uniqueIndex;not null"`                // Order ID unik (DONATION-123...)
	PaymentToken   string     `gorm:"type:text"`                           // Snap token
	PaymentGateway string     `gorm:"type:varchar(50);default:'midtrans'"` // Bisa juga xendit, dll
	PaymentMethod  string     `gorm:"type:varchar(50)"`                    // e.g. gopay, bca_va
	PaidAt         *time.Time // Waktu pembayaran sukses
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"` // Soft delete (opsional)
}

