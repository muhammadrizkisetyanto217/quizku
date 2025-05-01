package model

import "time"

type Tooltip struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Keyword          string    `gorm:"unique;not null" json:"keyword"`
	DescriptionShort string    `gorm:"type:varchar(200);not null" json:"description_short"`
	DescriptionLong  string    `gorm:"type:text;not null" json:"description_long"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName memastikan nama tabel sesuai dengan skema database
func (Tooltip) TableName() string {
	return "tooltips"
}