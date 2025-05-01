package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ReadingModel struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Title           string         `gorm:"type:varchar(50);unique;not null" json:"title"`
	Status          string         `gorm:"type:varchar(10);default:'pending';check:status IN ('active', 'pending', 'archived')" json:"status"`
	DescriptionLong string         `gorm:"type:text;not null" json:"description_long"`
	TooltipsID      pq.Int64Array  `gorm:"type:int[]" json:"tooltips_id"`
	CreatedAt       time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UnitID          uint           `json:"unit_id"`
	CreatedBy       uuid.UUID      `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"created_by"`
}

func (ReadingModel) TableName() string {
	return "readings"
}

func (r ReadingModel) MarshalJSONReading() ([]byte, error) {
	type Alias ReadingModel
	return json.Marshal(&struct {
		TooltipsID []int64 `json:"tooltips_id"`
		*Alias
	}{
		TooltipsID: []int64(r.TooltipsID), // 🔥 Konversi `pq.Int64Array` ke `[]int64`
		Alias:      (*Alias)(&r),
	})
}
