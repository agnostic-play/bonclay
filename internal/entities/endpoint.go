package entities

import (
	"strings"

	"gorm.io/gorm"
)

// EndpointEntity represents an endpoint in the system for API V2
type EndpointEntity struct {
	BaseEntityWithID

	Path           string `json:"path" gorm:"type:varchar(500);not null" validate:"omitempty,required"`
	Method         string `json:"method" gorm:"type:varchar(10);not null" validate:"omitempty,required,alpha"`
	CollectionID   string `json:"collection_id" gorm:"type:uuid;not null;index" validate:"omitempty,uuid"`
	Category       string `json:"category" gorm:"type:varchar(255)" validate:"omitempty"`
	ActiveScenario string `json:"active_scenario,omitempty" gorm:"column:active_scenario;type:uuid"`
	Desc           string `json:"desc,omitempty" gorm:"type:text"`
	Delay          *int   `json:"delay,omitempty" gorm:"column:delay"`

	EnableResponseIntercept bool `json:"enable_response_intercept" gorm:"column:enable_response_intercept;default:false"`
}

func (e *EndpointEntity) GetTableName() string {
	return "endpoints"
}

func (e *EndpointEntity) GetTitle() string {
	return "Endpoint"
}

func (e *EndpointEntity) GetPreloadTables() []string {
	return []string{}
}

func (e *EndpointEntity) GetFieldForKeywords() []string {
	return []string{"path", "method", "category", "desc"}
}

func (e *EndpointEntity) GetExcludeFieldForUpdate() []string {
	return []string{"collection_id", "active_scenario"}
}

func (e *EndpointEntity) GetEntity() *EndpointEntity {
	return e
}

func (e *EndpointEntity) BeforeCreate(tx *gorm.DB) error {

	if err := e.BaseEntityWithID.BeforeCreate(tx); err != nil {
		return err
	}
	e.Path = strings.ReplaceAll(e.Path, "{query_params}", "[^/]+")
	return nil
}
