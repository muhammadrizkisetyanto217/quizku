package service

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	model "quizku/internals/features/certificates/issued_certificates/model"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// func CreateIssuedCertificateIfEligible(
// 	db *gorm.DB,
// 	userID uuid.UUID,
// 	subcategoryID int,
// ) error {
// 	var certVersion struct {
// 		ID     uint
// 		Number int
// 	}
// 	err := db.Table("certificate_versions").
// 		Select("id, version_number").
// 		Where("subcategory_id = ?", subcategoryID).
// 		Order("version_number DESC").
// 		Limit(1).
// 		Scan(&certVersion).Error
// 	if err != nil {
// 		return err
// 	}

// 	var exists int64
// 	db.Table("issued_certificates").
// 		Where("user_id = ? AND subcategory_id = ?", userID, subcategoryID).
// 		Count(&exists)
// 	if exists > 0 {
// 		return nil
// 	}

// 	slug := fmt.Sprintf("cert-%s-%d", userID.String(), time.Now().Unix())
// 	issued := model.IssuedCertificateModel{
// 		UserID:               userID,
// 		SubcategoryID:        uint(subcategoryID),
// 		CertificateVersionID: certVersion.ID,
// 		VersionIssued:        certVersion.Number,
// 		VersionCurrent:       certVersion.Number,
// 		IsUpToDate:           true,
// 		SlugURL:              slug,
// 		IssuedAt:             time.Now(),
// 		CreatedAt:            time.Now(),
// 		UpdatedAt:            time.Now(),
// 	}

// 	if err := db.Create(&issued).Error; err != nil {
// 		log.Println("[ERROR] Gagal menyimpan issued certificate:", err)
// 		return err
// 	}

// 	return nil
// }
