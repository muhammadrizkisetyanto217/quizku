package model

import "time"

type CertificateVersionModel struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	SubcategoryID uint       `json:"subcategory_id"`
	VersionNumber int        `json:"version_number"`
	TotalThemes   int        `json:"total_themes"`
	Note          string     `json:"note"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

func (CertificateVersionModel) TableName() string {
	return "certificate_versions"
}
