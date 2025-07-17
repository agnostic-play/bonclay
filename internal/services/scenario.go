package services

import (
	"context"
	"github.com/agnostic-play/ditoo/internal/repository"
	"strings"
)

type ScenarioEntityReq struct {
	EndpointId   string `json:"endpoint_id" validate:"required,uuid"`
	Desc         string `json:"desc" validate:"required,customText" `
	Header       string `json:"header" validate:"required,customJSON"`
	Body         string `json:"body" validate:""`
	StatusHeader int    `json:"status_header" validate:"required,numeric"`
	Delay        *int   `json:"delay" validate:"required"`
}

func (req ScenarioEntityReq) translate() repository.ScenarioEntity {
	var ent repository.ScenarioEntity

	ent.Delay = req.Delay
	ent.Header = req.Header
	ent.Body = req.Body
	ent.StatusHeader = req.StatusHeader

	if val := strings.TrimSpace(req.Desc); val != "" {
		ent.Desc = strings.ToLower(val)
	}

	if val := strings.TrimSpace(req.EndpointId); val != "" {
		ent.EndpointId = val
	}
	return ent
}

type ScenarioServiceInterface interface {
	CreateOrUpdateScenario(ctx context.Context, squadID string, req ScenarioEntityReq) (repository.ScenarioEntity, error)
	DeleteScenario(ctx context.Context, squadID string) error
}

func (cont serviceContainer) CreateOrUpdateScenario(ctx context.Context, id string, req ScenarioEntityReq) (repository.ScenarioEntity, error) {
	resp, err := cont.repoContainer.CreateOrUpdateScenario(ctx, id, req.translate())
	if err != nil {
		return repository.ScenarioEntity{}, err
	}
	return resp, nil
}

func (cont serviceContainer) DeleteScenario(ctx context.Context, id string) error {
	_, err := cont.repoContainer.GetScenario(ctx, id)
	if err != nil {
		return err
	}

	err = cont.repoContainer.DeleteScenario(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
