package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/agnostic-play/ditoo/internal/repository"
)

type MockServiceInterface interface {
	MockApi(ctx context.Context, collectionSlug, method string, path string) (repository.ScenarioEntity, error)
}

func (cont serviceContainer) MockApi(ctx context.Context, collectionSlug, method string, path string) (repository.ScenarioEntity, error) {
	path = fmt.Sprintf("/%s", path)
	collection, err := cont.repoContainer.GetCollectionBySlug(ctx, collectionSlug)
	if err != nil {
		err = fmt.Errorf("collection not found")
		return repository.ScenarioEntity{}, err
	}

	endpoint, err := cont.repoContainer.GetEndpointMock(ctx, collection.ID.String(), method, path)
	if err != nil {
		return repository.ScenarioEntity{}, err
	}

	if strings.TrimSpace(endpoint.ActiveScenario) == "" {
		err = fmt.Errorf("scenario not found")
		return repository.ScenarioEntity{}, err
	}

	scenario, err := cont.repoContainer.GetScenario(ctx, endpoint.ActiveScenario)
	if err != nil {
		return repository.ScenarioEntity{}, err
	}

	customVar, err := cont.repoContainer.GetListCustomVariableByCollectionId(ctx, collection.ID.String())
	if err != nil {
		return repository.ScenarioEntity{}, err
	}
	scenario.ApplyEnv(envSliceToMap(customVar))

	return scenario, nil
}

func envSliceToMap(envs []repository.CustomVariableEntity) map[string]string {
	result := make(map[string]string)
	for _, e := range envs {
		result[e.Key] = e.Value
	}
	return result
}
