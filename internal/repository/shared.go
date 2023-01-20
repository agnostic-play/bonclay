package repository

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
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

func (cont repoContainerGorm) generateSlug(table, name string) string {
	slug := strings.Replace(strings.ToLower(strings.TrimSpace(name)), " ", "-", -1)

	var count int64
	if err := cont.db.Table(table).
		Where("name = ?", name).
		Count(&count).Error; err != nil {
		return slug
	}
	if count > 0 {
		slug = fmt.Sprintf("%s-%d", slug, count+1)
	}

	return slug
}
