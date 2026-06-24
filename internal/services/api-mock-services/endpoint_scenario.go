package api_mock_services

import (
	"context"
	"fmt"

	repository2 `berlin.allobank.com/tools/bonclay/internal/adapters/repositories`
	"berlin.allobank.com/tools/bonclay/internal/common/errs"
	"berlin.allobank.com/tools/bonclay/internal/common/pagination"
	"berlin.allobank.com/tools/bonclay/internal/entities"
	crud_services "berlin.allobank.com/tools/bonclay/internal/services/crud-services"
)

type EndpointScenarioService interface {
	GetEndpointScenario(ctx context.Context, slug string) (repository2.CollectionDetailResponse, error)
	SetActiveResponse(ctx context.Context, req SetActiveScenarioReq) error
	RemoveScenario(ctx context.Context, req SetActiveScenarioReq) error
	UpdateScript(ctx context.Context, endpointID, script string) error
}

type endpointScenarioService struct {
	collectionService crud_services.BaseCRUDService[entities.CollectionEntity]
	endpointService   crud_services.BaseCRUDService[entities.EndpointEntity]
	scenarioService   crud_services.BaseCRUDService[entities.ScenarioEntity]
	repoContainer     repository2.RepoContainer
}

func NewEndpointScenarioService(
	collectionService crud_services.BaseCRUDService[entities.CollectionEntity],
	endpointService crud_services.BaseCRUDService[entities.EndpointEntity],
	scenarioService crud_services.BaseCRUDService[entities.ScenarioEntity],
	repoContainer repository2.RepoContainer,
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

// UpdateScript validates the endpoint exists and persists its mock script.
func (s *endpointScenarioService) UpdateScript(ctx context.Context, endpointID, script string) error {
	if _, err := s.endpointService.Get(ctx, endpointID); err != nil {
		return err
	}

	if err := s.repoContainer.GetEndpointRepo().UpdateScript(ctx, endpointID, script); err != nil {
		return errs.BadRequest(fmt.Sprintf("can not update script: %s", err.Error()))
	}
	return nil
}

func (s *endpointScenarioService) GetEndpointScenario(ctx context.Context, slug string) (repository2.CollectionDetailResponse, error) {
	empty := repository2.CollectionDetailResponse{}

	// Try slug first; fall back to ID for collections that have no slug yet
	collQuery := &pagination.ListQuery{ShowAll: true, Filters: map[string]interface{}{"slug = ?": slug}}
	collResult, err := s.collectionService.GetList(ctx, collQuery)
	if err != nil || len(collResult.List) == 0 {
		collQuery = &pagination.ListQuery{ShowAll: true, Filters: map[string]interface{}{"id = ?": slug}}
		collResult, err = s.collectionService.GetList(ctx, collQuery)
		if err != nil || len(collResult.List) == 0 {
			return empty, fmt.Errorf("collection not found")
		}
	}
	collection := collResult.List[0]

	// Get endpoint scenarios using repository
	endpoints, err := s.repoContainer.GetEndpointScenarioRepo().GetListEndpointScenario(ctx, collection.ID.String())
	if err != nil {
		return empty, err
	}

	grouped := make(map[string][]repository2.EndpointScenario)
	for _, val := range endpoints {
		grouped[val.Category] = append(grouped[val.Category], val)
	}

	return repository2.CollectionDetailResponse{
		ID:              collection.ID.String(),
		Name:            collection.Name,
		Slug:            collection.Slug,
		Desc:            collection.Desc,
		Docs:            collection.Docs,
		ForwardProxyURL: collection.ForwardProxyURL,
		IsProxyEnable:   collection.IsProxyEnable,
		Endpoints:       grouped,
	}, nil
}
