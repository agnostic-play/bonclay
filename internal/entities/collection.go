package entities

import (
	"regexp"
	"strings"

	"gorm.io/gorm"
)

// CollectionEntity represents a collection in the system for API V2
type CollectionEntity struct {
	BaseEntityWithID

	Name    string `json:"name" gorm:"type:varchar(255);not null;" validate:"omitempty,required,min=3,max=255"`
	Slug    string `json:"slug" gorm:"type:varchar(255);uniqueIndex"`
	Docs    string `json:"docs,omitempty" gorm:"type:text"`
	SquadID string `json:"squad_id" gorm:"type:uuid;not null;index" validate:"omitempty,required,uuid"`
	Desc    string `json:"desc,omitempty" gorm:"type:text"`

	ForwardProxyURL string `json:"forward_proxy_url,omitempty" gorm:"column:forward_proxy_url;type:varchar(500)"`
	IsProxyEnable   *bool  `json:"is_proxy_enable" gorm:"column:is_proxy_enable;default:false"`
}

func (e *CollectionEntity) GetTableName() string {
	return "collections"
}

func (e *CollectionEntity) GetTitle() string {
	return "Collection"
}

func (e *CollectionEntity) GetPreloadTables() []string {
	return []string{}
}

func (e *CollectionEntity) GetFieldForKeywords() []string {
	return []string{"name", "desc", "docs"}
}

func (e *CollectionEntity) GetExcludeFieldForUpdate() []string {
	return []string{"slug", "squad_id"}
}

func (e *CollectionEntity) GetEntity() *CollectionEntity {
	return e
}

func (e *CollectionEntity) BeforeCreate(tx *gorm.DB) error {
	if err := e.BaseEntityWithID.BeforeCreate(tx); err != nil {
		return err
	}
	if e.Slug == "" {
		e.Slug = toSlug(e.Name)
	}
	return nil
}

var nonAlphanumRe = regexp.MustCompile(`[^a-z0-9]+`)

func toSlug(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = nonAlphanumRe.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
