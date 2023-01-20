package services

import (
	"gitlab.linkaja.com/be/ditto/internal/config"
	"gitlab.linkaja.com/be/ditto/internal/repository"
)

type ServiceContainer interface {
	SquadServiceInterface
	CollectionServiceInterface
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
