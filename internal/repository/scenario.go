package repository

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"gorm.io/gorm"
)

const scenarioResponse = "scenario_response"

type ScenarioEntity struct {
	BaseEntityWithID
	EndpointId   string `json:"endpoint_id" gorm:"column:endpoint_id"`
	Desc         string `json:"desc" gorm:"column:desc"`
	Header       string `json:"header" gorm:"column:header"`
	Body         string `json:"body" gorm:"column:body"`
	StatusHeader int    `json:"status_header" gorm:"column:status_header"`
	Delay        *int   `json:"delay" gorm:"column:delay"`
}

func (ScenarioEntity) TableName() string {
	return "scenario_response"
}

type ScenarioRepoInterface interface {
	GetScenario(ctx context.Context, id string) (ScenarioEntity, error)
	CreateOrUpdateScenario(ctx context.Context, squadID string, entity ScenarioEntity) (ScenarioEntity, error)
	DeleteScenario(ctx context.Context, squadID string) error
}

func (cont repoContainerGorm) CreateOrUpdateScenario(ctx context.Context, scenarioID string, entity ScenarioEntity) (ScenarioEntity, error) {
	exec := cont.write(ctx).Table(scenarioResponse)

	if scenarioID == "" {
		exec = exec.Create(&entity)
	} else {
		exec = exec.Omit("endpoint_id").Where("id=?", scenarioID).Updates(&entity)
	}

	if exec.Error != nil {
		return ScenarioEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	if exec.RowsAffected < 1 {
		return ScenarioEntity{}, fmt.Errorf("cannot create data")
	}

	return entity, nil
}

func (cont repoContainerGorm) GetListScenario(ctx context.Context, squadID string) ([]ScenarioEntity, error) {

	var entity []ScenarioEntity

	exec := cont.db.WithContext(ctx).Table(scenarioResponse).Where("squad_id = ?", squadID).Find(&entity)

	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return []ScenarioEntity{}, fmt.Errorf("collections not found")
		}
		return []ScenarioEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (cont repoContainerGorm) GetScenario(ctx context.Context, slug string) (ScenarioEntity, error) {

	var entity ScenarioEntity

	exec := cont.db.WithContext(ctx).Table(scenarioResponse).Where("id=?", slug).First(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return ScenarioEntity{}, fmt.Errorf("collections not found")
		}
		return ScenarioEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (cont repoContainerGorm) DeleteScenario(ctx context.Context, squadID string) error {

	var entity ScenarioEntity

	exec := cont.write(ctx).Table(scenarioResponse).Where("id=?", squadID).Delete(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("collections not found")
		}
		return fmt.Errorf("query error: %s", exec.Error)
	}

	return nil
}

func (scenario *ScenarioEntity) ApplyEnv(env map[string]string) {
	// Match {{varName}}
	re := regexp.MustCompile(`\{\{(.+?)\}\}`)

	// Replace each match with env value
	result := re.ReplaceAllStringFunc(scenario.Body, func(m string) string {
		// Extract key name without {{ }}
		key := re.FindStringSubmatch(m)[1]
		if val, ok := env[key]; ok {
			return val
		}
		// If not found, keep original placeholder
		return m
	})

	scenario.Body = result

	return
}
