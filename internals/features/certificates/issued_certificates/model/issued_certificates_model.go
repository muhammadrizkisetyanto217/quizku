package model

import (
	"time"

	"github.com/google/uuid"
)

type IssuedCertificateModel struct {
	ID                  uint      `json:"id" gorm:"primaryKey"`
	UserID              uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	SubcategoryID       uint      `json:"subcategory_id" gorm:"not null"`
	CertificateVersionID uint     `json:"certificate_version_id" gorm:"not null"`
	VersionCurrent      int       `json:"version_current" gorm:"not null"`
	IsUpToDate          bool      `json:"is_up_to_date" gorm:"not null;default:true"`
	SlugURL             string    `json:"slug_url" gorm:"unique;not null"`
	IssuedAt            time.Time `json:"issued_at" gorm:"not null"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (IssuedCertificateModel) TableName() string {
	return "issued_certificates"
}
