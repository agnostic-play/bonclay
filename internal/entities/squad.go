package entities

// SquadEntity represents a squad in the system for API V2
type SquadEntity struct {
	BaseEntityWithID

	Name     string `json:"name" gorm:"type:varchar(255);not null;" validate:"required,min=3,max=255"`
	Slug     string `json:"slug" gorm:"type:varchar(255);uniqueIndex"`
	Password string `json:"-" gorm:"type:varchar(255)"`
	Desc     string `json:"desc,omitempty" gorm:"type:text"`
}

func (e SquadEntity) GetTableName() string {
	return "squads"
}

func (e SquadEntity) GetTitle() string {
	return "Squad"
}

func (e SquadEntity) GetPreloadTables() []string {
	return []string{}
}

func (e SquadEntity) GetFieldForKeywords() []string {
	return []string{"name", "desc"}
}

func (e SquadEntity) GetExcludeFieldForUpdate() []string {
	return []string{"slug", "password"}
}

func (e SquadEntity) GetEntity() *SquadEntity {
	return &e
}
