package entities

// EndpointEntity represents an endpoint in the system for API V2
type EndpointEntity struct {
	BaseEntityWithID

	Path           string `json:"path" gorm:"type:varchar(500);not null" validate:"required"`
	Method         string `json:"method" gorm:"type:varchar(10);not null" validate:"required,alpha"`
	CollectionID   string `json:"collection_id" gorm:"type:uuid;not null;index" validate:"required,uuid"`
	Category       string `json:"category" gorm:"type:varchar(255)" validate:"omitempty"`
	ActiveScenario string `json:"active_scenario,omitempty" gorm:"column:active_scenario;type:uuid"`
	Desc           string `json:"desc,omitempty" gorm:"type:text"`

	EnableResponseIntercept bool `json:"enable_response_intercept" gorm:"column:enable_response_intercept;default:false"`
}

func (e EndpointEntity) GetTableName() string {
	return "endpoints"
}

func (e EndpointEntity) GetTitle() string {
	return "Endpoint"
}

func (e EndpointEntity) GetPreloadTables() []string {
	return []string{}
}

func (e EndpointEntity) GetFieldForKeywords() []string {
	return []string{"path", "method", "category", "desc"}
}

func (e EndpointEntity) GetExcludeFieldForUpdate() []string {
	return []string{"collection_id", "active_scenario"}
}

func (e EndpointEntity) GetEntity() *EndpointEntity {
	return &e
}
