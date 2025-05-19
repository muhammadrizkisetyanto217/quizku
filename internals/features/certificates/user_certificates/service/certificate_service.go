package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateOrUpdateIssuedCertificate(
	db *gorm.DB,
	userID uuid.UUID,
	subcategoryID int,
	issuedVersion int, // tetap diterima untuk keperluan eksternal
) error {
	var existingID uint
	err := db.Table("issued_certificates").
		Select("id").
		Where("user_id = ? AND subcategory_id = ?", userID, subcategoryID).
		Scan(&existingID).Error
	if err != nil {
		return err
	}

	now := time.Now()

	if existingID > 0 {
		// Update existing (tanpa version info)
		return db.Table("issued_certificates").
			Where("id = ?", existingID).
			Updates(map[string]interface{}{
				"user_cert_is_up_to_date": true,
				"certificate_issued_at":   now,
				"updated_at":              now,
			}).Error
	}

	// Create new
	slug := fmt.Sprintf("cert-%s-%d", userID.String(), now.Unix())
	return db.Table("issued_certificates").Create(map[string]interface{}{
		"user_id":                 userID,
		"subcategory_id":          subcategoryID,
		"user_cert_is_up_to_date": true,
		"certificate_slug_url":    slug,
		"certificate_issued_at":   now,
		"created_at":              now,
		"updated_at":              now,
	}).Error
}