package repositories

import (
	"berlin.allobank.com/tools/bonclay/internal/config"
	"gorm.io/gorm"
)

// RepoContainer - aggregates all repository interfaces
type RepoContainer interface {
	GetEndpointRepo() EndpointRepoInterface
	GetEndpointScenarioRepo() EndpointScenarioRepoInterface
	GetCustomVariableRepo() CustomVariableRepoInterface
}

type repoContainer struct {
	endpointRepo         EndpointRepoInterface
	endpointScenarioRepo EndpointScenarioRepoInterface
	customVariableRepo   CustomVariableRepoInterface
}

func NewRepoContainer(db *gorm.DB, config *config.Config) RepoContainer {
	dbClient := NewDBClient(db)

	return &repoContainer{
		endpointRepo:         NewEndpointRepo(dbClient),
		endpointScenarioRepo: NewEndpointScenarioRepo(dbClient),
		customVariableRepo:   NewCustomVariableRepo(dbClient),
	}
}

func (c *repoContainer) GetEndpointRepo() EndpointRepoInterface {
	return c.endpointRepo
}

func (c *repoContainer) GetEndpointScenarioRepo() EndpointScenarioRepoInterface {
	return c.endpointScenarioRepo
}

func (c *repoContainer) GetCustomVariableRepo() CustomVariableRepoInterface {
	return c.customVariableRepo
}
