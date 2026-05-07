package entities

import (
	"gorm.io/gorm"
)

// CustomVariable represents an endpoint in the system for API V2
type CustomVariable struct {
	BaseEntityWithID
	CollectionID string `json:"collection_id" gorm:"column:collection_id;unique" validate:"omitempty,required,uuid"`
	Key          string `json:"key" gorm:"column:key;unique" validate:"omitempty,required,uuid"`
	Value        string `json:"value" gorm:"column:value;unique" validate:"omitempty,required,uuid"`
	IsActive     bool   `json:"is_active" gorm:"column:is_active"`
}

func (e *CustomVariable) GetTableName() string {
	return "custom_variables"
}

func (e *CustomVariable) GetTitle() string {
	return "Custom Variable"
}

func (e *CustomVariable) GetPreloadTables() []string {
	return []string{}
}

func (e *CustomVariable) GetFieldForKeywords() []string {
	return []string{"key"}
}

func (e *CustomVariable) GetExcludeFieldForUpdate() []string {
	return []string{"key"}
}

func (e *CustomVariable) GetEntity() *CustomVariable {
	return e
}

func (e *CustomVariable) BeforeCreate(tx *gorm.DB) error {
	if err := e.BaseEntityWithID.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

func (e *CustomVariable) BeforeUpdate(tx *gorm.DB) (err error) {
	if err := e.BaseEntityWithID.BeforeUpdate(tx); err != nil {
		return err
	}

	return nil
}
