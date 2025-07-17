package services

import (
	"context"
	"fmt"
	"github.com/agnostic-play/ditoo/internal/repository"
	"strings"
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

	return scenario, nil
}
