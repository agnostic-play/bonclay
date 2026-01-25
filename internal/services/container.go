package services

import (
	"github.com/agnostic-play/ditoo/internal/config"
	"github.com/agnostic-play/ditoo/internal/repository"
	api_mock_services "github.com/agnostic-play/ditoo/internal/services/api-mock-services"
	common_services "github.com/agnostic-play/ditoo/internal/services/common-services"
	crud_services "github.com/agnostic-play/ditoo/internal/services/crud-services"
)

// ServiceContainer provides access to all application services
type ServiceContainer interface {
	// CRUD Services
	GetCRUDServices() *crud_services.Services

	// API Mock Services
	GetMockService() api_mock_services.MockService
	GetEndpointScenarioService() api_mock_services.EndpointScenarioService
	GetCustomVariableService() api_mock_services.CustomVariableService

	// Common Services
	GetEncryptionService() common_services.EncryptionService
}

func NewServiceContainer(
	cfg *config.Config,
	repoContainer repository.RepoContainer,
	dbClient repository.DBClient,
) ServiceContainer {
	// Create CRUD services first since they're used by api-mock-services
	crudServices := crud_services.NewCRUDServices(dbClient)

	return &serviceContainer{
		crudServices: crudServices,
		mockService: api_mock_services.NewMockService(
			crudServices.CollectionServices,
			crudServices.ScenarioServices,
			repoContainer,
		),
		endpointScenarioService: api_mock_services.NewEndpointScenarioService(
			crudServices.CollectionServices,
			crudServices.EndpointServices,
			crudServices.ScenarioServices,
			repoContainer,
		),
		customVariableService: api_mock_services.NewCustomVariableService(
			crudServices.CollectionServices,
			repoContainer,
		),
		encryptionService: common_services.NewEncryptionService(),
		config:            cfg,
	}
}

type serviceContainer struct {
	crudServices            *crud_services.Services
	mockService             api_mock_services.MockService
	endpointScenarioService api_mock_services.EndpointScenarioService
	customVariableService   api_mock_services.CustomVariableService
	encryptionService       common_services.EncryptionService
	config                  *config.Config
}

func (c *serviceContainer) GetCRUDServices() *crud_services.Services {
	return c.crudServices
}

func (c *serviceContainer) GetMockService() api_mock_services.MockService {
	return c.mockService
}

func (c *serviceContainer) GetEndpointScenarioService() api_mock_services.EndpointScenarioService {
	return c.endpointScenarioService
}

func (c *serviceContainer) GetCustomVariableService() api_mock_services.CustomVariableService {
	return c.customVariableService
}

func (c *serviceContainer) GetEncryptionService() common_services.EncryptionService {
	return c.encryptionService
}
