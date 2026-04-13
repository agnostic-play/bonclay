package crud_services

import (
	`github.com/agnostic-play/ditoo/internal/adapters/repositories`
	"github.com/agnostic-play/ditoo/internal/entities"
)

type Services struct {
	DiagramCollectionServices BaseCRUDService[entities.DiagramCollectionEntity]
	DiagramServices           BaseCRUDService[entities.DiagramEntity]

	// V2 API Services
	SquadServices      BaseCRUDService[entities.SquadEntity]
	CollectionServices BaseCRUDService[entities.CollectionEntity]
	EndpointServices   BaseCRUDService[entities.EndpointEntity]
	ScenarioServices   BaseCRUDService[entities.ScenarioEntity]
}

func NewCRUDServices(dbClient repositories.DBClient) *Services {
	return &Services{
		DiagramCollectionServices: NewBaseCRUDServices[entities.DiagramCollectionEntity](new(entities.DiagramCollectionEntity), dbClient),
		DiagramServices:           NewBaseCRUDServices[entities.DiagramEntity](new(entities.DiagramEntity), dbClient),

		// V2 API Services
		SquadServices:      NewBaseCRUDServices[entities.SquadEntity](new(entities.SquadEntity), dbClient),
		CollectionServices: NewBaseCRUDServices[entities.CollectionEntity](new(entities.CollectionEntity), dbClient),
		EndpointServices:   NewBaseCRUDServices[entities.EndpointEntity](new(entities.EndpointEntity), dbClient),
		ScenarioServices:   NewBaseCRUDServices[entities.ScenarioEntity](new(entities.ScenarioEntity), dbClient),
	}
}
