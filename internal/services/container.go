package services

import (
	"github.com/agnostic-play/ditoo/internal/config"
	"github.com/agnostic-play/ditoo/internal/repository"
	crud_services "github.com/agnostic-play/ditoo/internal/services/crud-services"
)

type ServiceContainer interface {
	GetCRUDServices() *crud_services.Services
	SquadServiceInterface
	CollectionServiceInterface
	EndpointServiceInterface
	ScenarioServiceInterface
	EndpointScenarioServiceInterface
	MockServiceInterface
	ToolsService
}

func NewServiceContainer(config *config.Config,
	repoContainer repository.RepoContainer,
	dbClient repository.DBClient) ServiceContainer {
	
	return &serviceContainer{
		repoContainer: repoContainer,
		crudServices:  crud_services.NewCRUDServices(dbClient),
		config:        config,
	}
}

type serviceContainer struct {
	repoContainer repository.RepoContainer
	crudServices  *crud_services.Services
	config        *config.Config
}

func (c serviceContainer) GetCRUDServices() *crud_services.Services {
	return c.crudServices
}
