package services

import (
	"github.com/agnostic-play/ditoo/internal/config"
	"github.com/agnostic-play/ditoo/internal/repository"
)

type ServiceContainer interface {
	SquadServiceInterface
	CollectionServiceInterface
	EndpointServiceInterface
	ScenarioServiceInterface
	EndpointScenarioServiceInterface
	CustomVariableServiceInterface
	MockServiceInterface
	ToolsService
}

func NewServiceContainer(repoContainer repository.RepoContainer,
	config *config.Config) ServiceContainer {
	return &serviceContainer{
		repoContainer: repoContainer,
		config:        config,
	}
}

type serviceContainer struct {
	repoContainer repository.RepoContainer
	config        *config.Config
}
