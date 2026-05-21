package services

import (
	"fmt"

	redisclient "berlin.allobank.com/tools/bonclay/internal/adapters/redis"
	repository2 "berlin.allobank.com/tools/bonclay/internal/adapters/repositories"
	"berlin.allobank.com/tools/bonclay/internal/config"
	api_mock_services "berlin.allobank.com/tools/bonclay/internal/services/api-mock-services"
	common_services "berlin.allobank.com/tools/bonclay/internal/services/common-services"
	crud_services "berlin.allobank.com/tools/bonclay/internal/services/crud-services"
	kong_services "berlin.allobank.com/tools/bonclay/internal/services/kong-services"
)

// ServiceContainer provides access to all application services
type ServiceContainer interface {
	// CRUD Services
	GetCRUDServices() *crud_services.Services

	// API Mock Services
	GetMockService() api_mock_services.MockService
	GetEndpointScenarioService() api_mock_services.EndpointScenarioService
	GetCustomVariableService() api_mock_services.CustomVariableService
	GetHistoryService() api_mock_services.HistoryService

	// Common Services
	GetEncryptionService() common_services.EncryptionService

	// Kong Services
	GetKongService() kong_services.KongService
}

func NewServiceContainer(
	cfg *config.Config,
	repoContainer repository2.RepoContainer,
	dbClient repository2.DBClient,
) ServiceContainer {
	crudServices := crud_services.NewCRUDServices(dbClient)

	redis := redisclient.NewRedisClient(
		fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)

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
		historyService:    api_mock_services.NewHistoryService(redis),
		encryptionService: common_services.NewEncryptionService(),
		kongService:       kong_services.NewKongService(cfg.Kong),
		config:            cfg,
	}
}

type serviceContainer struct {
	crudServices            *crud_services.Services
	mockService             api_mock_services.MockService
	endpointScenarioService api_mock_services.EndpointScenarioService
	customVariableService   api_mock_services.CustomVariableService
	historyService          api_mock_services.HistoryService
	encryptionService       common_services.EncryptionService
	kongService             kong_services.KongService
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

func (c *serviceContainer) GetHistoryService() api_mock_services.HistoryService {
	return c.historyService
}

func (c *serviceContainer) GetEncryptionService() common_services.EncryptionService {
	return c.encryptionService
}

func (c *serviceContainer) GetKongService() kong_services.KongService {
	return c.kongService
}
