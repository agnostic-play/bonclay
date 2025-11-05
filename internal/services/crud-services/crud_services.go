package crud_services

import (
	"github.com/agnostic-play/ditoo/internal/entities"
	"github.com/agnostic-play/ditoo/internal/repository"
)

type Services struct {
	DiagramCollectionServices BaseCRUDService[entities.DiagramCollectionEntity]
	DiagramServices           BaseCRUDService[entities.DiagramEntity]
}

func NewCRUDServices(dbClient repository.DBClient) *Services {
	return &Services{
		DiagramCollectionServices: NewBaseCRUDServices(new(entities.DiagramCollectionEntity), dbClient),
		DiagramServices:           NewBaseCRUDServices(new(entities.DiagramEntity), dbClient),
	}
}
