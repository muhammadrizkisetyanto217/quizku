package service

import (
	"quizku/internals/features/certificates/user_certificates/model"
	subcategoryModel "quizku/internals/features/lessons/subcategories/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	completed := len(us.UserSubcategoryCompleteThemesOrLevels)
	total := len(sub.SubcategoryTotalThemesOrLevels)

	// 🧠 Logika validasi apakah up-to-date
	isUpToDate := (us.UserSubcategoryCurrentVersion == issuedVersion) && (completed >= total)

	// 🔁 Update hanya jika status berubah
	if cert.UserCertIsUpToDate != isUpToDate {
		err := db.Model(&model.UserCertificate{}).
			Where("user_cert_id = ?", cert.UserCertID).
			Update("user_cert_is_up_to_date", isUpToDate).Error
		if err != nil {
			return cert.UserCertIsUpToDate, err
		}
	}

	return isUpToDate, nil
}
