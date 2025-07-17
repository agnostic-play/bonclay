package services

import (
	"context"
	"fmt"
	"github.com/agnostic-play/ditoo/internal/errs"
	"github.com/agnostic-play/ditoo/internal/repository"
)

type SetActiveScenarioEntityReq struct {
	EndpointID string `json:"endpoint_id" validate:"required,uuid"`
	ScenarioID string `json:"scenario_id" validate:"required,uuid"`
}

type EndpointScenarioServiceInterface interface {
	GetEndpointScenario(ctx context.Context, slug string) (map[string][]repository.EndpointScenario, error)
	SetActiveResponse(ctx context.Context, req SetActiveScenarioEntityReq) error
}

func (cont serviceContainer) SetActiveResponse(ctx context.Context, req SetActiveScenarioEntityReq) error {
	_, err := cont.repoContainer.GetEndpoint(ctx, req.EndpointID)
	if err != nil {
		return err
	}
	_, err = cont.repoContainer.GetScenario(ctx, req.ScenarioID)
	if err != nil {
		return err
	}

	err = cont.repoContainer.SetActiveResponse(ctx, req.EndpointID, req.ScenarioID)
	if err != nil {
		err := errs.BadRequest(fmt.Sprintf("can not set active response: %s", err.Error()))
		return err
	}
	return nil
}

func (cont serviceContainer) GetEndpointScenario(ctx context.Context, slug string) (map[string][]repository.EndpointScenario, error) {

	resp := make(map[string][]repository.EndpointScenario)

	collection, err := cont.repoContainer.GetCollectionBySlug(ctx, slug)
	if err != nil {
		return map[string][]repository.EndpointScenario{}, err
	}

	endpoints, err := cont.repoContainer.GetListEndpointScenario(ctx, collection.ID.String())
	if err != nil {
		return map[string][]repository.EndpointScenario{}, err
	}

	for _, val := range endpoints {
		if _, ok := resp[val.Category]; !ok {
			resp[val.Category] = []repository.EndpointScenario{}
		}
		resp[val.Category] = append(resp[val.Category], val)
	}

	return resp, nil
}
