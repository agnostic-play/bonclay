package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type EndpointScenario struct {
	EndpointEntity
	ScenarioResponse []ScenarioEntity `json:"scenario_response" gorm:"foreignKey:endpoint_id"`
}

type EndpointScenarioRepoInterface interface {
	GetListEndpointScenario(ctx context.Context, collectionID string) ([]EndpointScenario, error)
	SetActiveResponse(ctx context.Context, endpointID, scenarioID string) error
}

func (cont repoContainerGorm) SetActiveResponse(ctx context.Context, endpointID, scenarioID string) error {

	exec := cont.write(ctx).Table(endpoints).
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

func (cont repoContainerGorm) GetListEndpointScenario(ctx context.Context, collectionID string) ([]EndpointScenario, error) {
	var entity []EndpointScenario

	exec := cont.db.Model(&EndpointScenario{}).
		Where("collection_id = ?", collectionID).
		Preload("ScenarioResponse").Find(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return []EndpointScenario{}, fmt.Errorf("collections not found")
		}
		return []EndpointScenario{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}
