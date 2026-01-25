package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntityWithID struct {
	ID        uuid.UUID      `json:"id"  gorm:"column:id;type:uuid;primary_key"`
	CreatedAt time.Time      `json:"-"  gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"-"  gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" `
}

func (entity *BaseEntityWithID) BeforeCreate(tx *gorm.DB) (err error) {
	entity.ID = uuid.New()
	return
}
