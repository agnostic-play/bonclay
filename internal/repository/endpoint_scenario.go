package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/agnostic-play/ditoo/internal/entities"
	"gorm.io/gorm"
)

// EndpointScenario - join struct for endpoint with scenarios
type EndpointScenario struct {
	entities.EndpointEntity
	ScenarioResponse []entities.ScenarioEntity `json:"scenario_response" gorm:"foreignKey:endpoint_id"`
}

func (EndpointScenario) TableName() string {
	return "endpoints"
}

// EndpointScenarioRepoInterface - custom operations for endpoint-scenario management
type EndpointScenarioRepoInterface interface {
	GetListEndpointScenario(ctx context.Context, collectionID string) ([]EndpointScenario, error)
	SetActiveResponse(ctx context.Context, endpointID, scenarioID string) error
	RemoveScenario(ctx context.Context, endpointID string) error
}

type endpointScenarioRepository struct {
	dbClient DBClient
}

func NewEndpointScenarioRepo(dbClient DBClient) EndpointScenarioRepoInterface {
	return &endpointScenarioRepository{dbClient: dbClient}
}

func (r *endpointScenarioRepository) SetActiveResponse(ctx context.Context, endpointID, scenarioID string) error {
	exec := r.dbClient.client(ctx).Table("endpoints").
		Where("id = ?", endpointID).
		Update("active_scenario", scenarioID)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("endpoint not found")
		}
		return fmt.Errorf("query error: %s", exec.Error)
	}

	return nil
}

func (r *endpointScenarioRepository) RemoveScenario(ctx context.Context, endpointID string) error {
	exec := r.dbClient.client(ctx).Table("endpoints").
		Where("id = ?", endpointID).
		Update("active_scenario", "")
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("endpoint not found")
		}
		return fmt.Errorf("query error: %s", exec.Error)
	}

	return nil
}

func (r *endpointScenarioRepository) GetListEndpointScenario(ctx context.Context, collectionID string) ([]EndpointScenario, error) {
	var entity []EndpointScenario

	exec := r.dbClient.client(ctx).Model(&EndpointScenario{}).
		Where("collection_id = ?", collectionID).
		Preload("ScenarioResponse").Find(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return []EndpointScenario{}, fmt.Errorf("endpoint not found")
		}
		return []EndpointScenario{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}
