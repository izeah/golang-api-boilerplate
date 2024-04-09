package abstraction

import (
	"time"

	"boilerplate/pkg/date"

	"gorm.io/gorm"
)

type Entity struct {
	ID           int        `json:"id" param:"id" form:"id" validate:"number,min=1" gorm:"primaryKey;autoIncrement;"`
	CreatedDate  time.Time  `json:"created_date" gorm:"<-:create" example:"1945-08-17T10:00:00Z"`
	CreatedBy    int        `json:"created_by" gorm:"<-:create" example:"1"`
	ModifiedDate *time.Time `json:"modified_date" gorm:"<-:update" example:"1945-08-17T10:00:00Z"`
	ModifiedBy   *int       `json:"modified_by" gorm:"<-:update" example:"1"`
}

// BeforeUpdate ...
func (m *Entity) BeforeUpdate(tx *gorm.DB) (err error) {
	if m.ModifiedDate == nil {
		m.ModifiedDate = date.NowUTC()
	}
	return
}

// BeforeCreate ...
func (m *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	if m.CreatedDate.IsZero() {
		m.CreatedDate = *date.NowUTC()
	}
	return
}
