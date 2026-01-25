package api_mock_services

import (
	"context"
	"fmt"

	"github.com/agnostic-play/ditoo/internal/common/errs"
	"github.com/agnostic-play/ditoo/internal/common/pagination"
	"github.com/agnostic-play/ditoo/internal/entities"
	"github.com/agnostic-play/ditoo/internal/repository"
	crud_services "github.com/agnostic-play/ditoo/internal/services/crud-services"
)

type EndpointScenarioService interface {
	GetEndpointScenario(ctx context.Context, slug string) (map[string][]repository.EndpointScenario, error)
	SetActiveResponse(ctx context.Context, req SetActiveScenarioReq) error
	RemoveScenario(ctx context.Context, req SetActiveScenarioReq) error
}

type endpointScenarioService struct {
	collectionService crud_services.BaseCRUDService[entities.CollectionEntity]
	endpointService   crud_services.BaseCRUDService[entities.EndpointEntity]
	scenarioService   crud_services.BaseCRUDService[entities.ScenarioEntity]
	repoContainer     repository.RepoContainer
}

func NewEndpointScenarioService(
	collectionService crud_services.BaseCRUDService[entities.CollectionEntity],
	endpointService crud_services.BaseCRUDService[entities.EndpointEntity],
	scenarioService crud_services.BaseCRUDService[entities.ScenarioEntity],
	repoContainer repository.RepoContainer,
) EndpointScenarioService {
	return &endpointScenarioService{
		collectionService: collectionService,
		endpointService:   endpointService,
		scenarioService:   scenarioService,
		repoContainer:     repoContainer,
	}
}

func (s *endpointScenarioService) SetActiveResponse(ctx context.Context, req SetActiveScenarioReq) error {
	// Validate endpoint exists using CRUD service
	_, err := s.endpointService.Get(ctx, req.EndpointID)
	if err != nil {
		return err
	}

	// Validate scenario exists using CRUD service
	_, err = s.scenarioService.Get(ctx, req.ScenarioID)
	if err != nil {
		return err
	}

	// Set active response using repository
	err = s.repoContainer.GetEndpointScenarioRepo().SetActiveResponse(ctx, req.EndpointID, req.ScenarioID)
	if err != nil {
		return errs.BadRequest(fmt.Sprintf("can not set active response: %s", err.Error()))
	}
	return nil
}

func (s *endpointScenarioService) RemoveScenario(ctx context.Context, req SetActiveScenarioReq) error {
	// Validate endpoint exists using CRUD service
	_, err := s.endpointService.Get(ctx, req.EndpointID)
	if err != nil {
		return err
	}

	// Remove scenario using repository
	err = s.repoContainer.GetEndpointScenarioRepo().RemoveScenario(ctx, req.EndpointID)
	if err != nil {
		return errs.BadRequest(fmt.Sprintf("can not remove active response: %s", err.Error()))
	}
	return nil
}

func (s *endpointScenarioService) GetEndpointScenario(ctx context.Context, slug string) (map[string][]repository.EndpointScenario, error) {
	resp := make(map[string][]repository.EndpointScenario)

	// Get collection by slug using CRUD service
	collQuery := &pagination.ListQuery{ShowAll: true, Filters: map[string]interface{}{"slug = ?": slug}}
	collResult, err := s.collectionService.GetList(ctx, collQuery)
	if err != nil || len(collResult.List) == 0 {
		return map[string][]repository.EndpointScenario{}, fmt.Errorf("collection not found")
	}
	collection := collResult.List[0]

	// Get endpoint scenarios using repository
	endpoints, err := s.repoContainer.GetEndpointScenarioRepo().GetListEndpointScenario(ctx, collection.ID.String())
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
