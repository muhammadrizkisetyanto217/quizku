package service

import (
	"quizku/internals/features/certificates/issued_certificates/model"

	"github.com/google/uuid"
	"gorm.io/gorm"

	subcategoryModel "quizku/internals/features/lessons/subcategories/model"
)


func CheckAndUpdateIsUpToDate(
	db *gorm.DB,
	userID uuid.UUID,
	subcategoryID int,
	cert model.IssuedCertificateModel,
	us subcategoryModel.UserSubcategoryModel,
	sub subcategoryModel.SubcategoryModel,
	issuedVersion int,
) (bool, error) {
	completed := len(us.CompleteThemesOrLevels)
	total := len(sub.TotalThemesOrLevels)

	// Logika validasi
	isUpToDate := (us.CurrentVersion == issuedVersion) && (completed >= total)

	// Hanya update jika nilai berbeda
	if cert.IsUpToDate != isUpToDate {
		err := db.Model(&model.IssuedCertificateModel{}).
			Where("id = ?", cert.ID).
			Update("is_up_to_date", isUpToDate).Error
		if err != nil {
			return cert.IsUpToDate, err
		}
	}

	return isUpToDate, nil
}
