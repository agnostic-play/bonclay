package entities

import (
	"github.com/google/uuid"
)

type DiagramCollectionEntity struct {
	BaseEntityWithID

	Name        string `json:"name" gorm:"type:varchar(255);not null" validate:"required,min=3,max=255" example:"System Architecture"`
	Description string `json:"description,omitempty" gorm:"type:text" validate:"omitempty,max=1000" example:"High-level architecture diagrams for backend systems."`
}

func (c DiagramCollectionEntity) GetTableName() string {
	return "diagram_collections"
}

func (c DiagramCollectionEntity) GetTitle() string {
	return "Diagram Collection"
}

func (c DiagramCollectionEntity) GetPreloadTables() []string {
	return []string{"Diagrams"}
}

func (c DiagramCollectionEntity) GetFieldForKeywords() []string {
	return []string{"name", "description"}
}

func (c DiagramCollectionEntity) GetEntity() *DiagramCollectionEntity {
	return &c
}

type DiagramEntity struct {
	BaseEntityWithID

	CollectionID uuid.UUID `json:"collection_id" gorm:"type:uuid;not null;index" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title        string    `json:"title" gorm:"type:varchar(255);not null" validate:"required,min=3,max=255" example:"Payment Flow Sequence Diagram"`
	Description  string    `json:"description,omitempty" gorm:"type:text" validate:"omitempty,max=1000" example:"Shows sequence of events between payment services."`
	SyntaxType   string    `json:"syntax_type" gorm:"type:varchar(50);default:'mermaid'" validate:"required,oneof=mermaid plantuml bpmn uml other" example:"mermaid"`
	Syntax       string    `json:"syntax" gorm:"type:text;not null" validate:"required,min=5" example:"graph TD; A-->B; B-->C;"`

	// Relations
	Collection *DiagramCollectionEntity `json:"collection,omitempty" gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE;"`
}

func (d DiagramEntity) GetTableName() string {
	return "diagrams"
}

func (d DiagramEntity) GetTitle() string {
	return "Diagram"
}

func (d DiagramEntity) GetPreloadTables() []string {
	return []string{"Collection"}
}

func (d DiagramEntity) GetFieldForKeywords() []string {
	return []string{"title", "description", "syntax_type"}
}

func (d DiagramEntity) GetEntity() *DiagramEntity {
	return &d
}
