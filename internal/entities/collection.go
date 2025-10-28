package entities

type CollectionEntity struct {
	BaseEntityWithID
	Name    string `json:"name"`
	Docs    string `json:"docs"`
	Slug    string `json:"slug"`
	SquadID string `json:"squad_id"`
	Desc    string `json:"desc"`
}

func (c CollectionEntity) GetTableName() string {
	return "collections"
}

func (c CollectionEntity) GetTitle() string {
	return "Collection"
}

func (c CollectionEntity) GetPreloadTables() []string {
	return []string{"collections"}
}

func (c CollectionEntity) GetFieldForKeywords() []string {
	return []string{"collections"}
}

func (c CollectionEntity) GetEntity() *CollectionEntity {
	return &c
}
