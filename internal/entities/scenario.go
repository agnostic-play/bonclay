package entities

// ScenarioEntity represents a scenario response in the system for API V2
type ScenarioEntity struct {
	BaseEntityWithID

	EndpointID   string `json:"endpoint_id" gorm:"column:endpoint_id;type:uuid;not null;index" validate:"required,uuid"`
	Desc         string `json:"desc" gorm:"type:text" validate:"required"`
	Header       string `json:"header" gorm:"type:text" validate:"required"`
	Body         string `json:"body" gorm:"type:text"`
	StatusHeader int    `json:"status_header" gorm:"column:status_header" validate:"required,numeric"`
	Delay        *int   `json:"delay" gorm:"column:delay"`
}

func (e ScenarioEntity) GetTableName() string {
	return "scenario_response"
}

func (e ScenarioEntity) GetTitle() string {
	return "Scenario"
}

func (e ScenarioEntity) GetPreloadTables() []string {
	return []string{}
}

func (e ScenarioEntity) GetFieldForKeywords() []string {
	return []string{"desc"}
}

func (e ScenarioEntity) GetExcludeFieldForUpdate() []string {
	return []string{"endpoint_id"}
}

func (e ScenarioEntity) GetEntity() *ScenarioEntity {
	return &e
}
