package service

import (
	"quizku/internals/features/certificates/user_certificates/model"

	"github.com/google/uuid"
	"gorm.io/gorm"

	subcategoryModel "quizku/internals/features/lessons/subcategories/model"
)

func CheckAndUpdateIsUpToDate(
	db *gorm.DB,
	userID uuid.UUID,
	subcategoryID int,
	cert model.UserCertificate,
	us subcategoryModel.UserSubcategoryModel,
	sub subcategoryModel.SubcategoryModel,
	issuedVersion int,
) (bool, error) {
	completed := len(us.CompleteThemesOrLevels)
	total := len(sub.TotalThemesOrLevels)

	// Logika validasi
	isUpToDate := (us.CurrentVersion == issuedVersion) && (completed >= total)

	// Hanya update jika nilai berbeda
	if cert.UserCertIsUpToDate != isUpToDate {
		err := db.Model(&model.UserCertificate{}).
			Where("id = ?", cert.UserCertID).
			Update("user_cert_is_up_to_date", isUpToDate).Error
		if err != nil {
			return cert.UserCertIsUpToDate, err
		}
	}

	return isUpToDate, nil
}
