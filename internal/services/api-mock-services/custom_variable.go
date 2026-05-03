package api_mock_services

import (
	"context"
	"fmt"
	"net/http"

	repository2 `berlin.allobank.com/tools/bonclay/internal/adapters/repositories`
	"berlin.allobank.com/tools/bonclay/internal/common/errs"
	"berlin.allobank.com/tools/bonclay/internal/common/pagination"
	"berlin.allobank.com/tools/bonclay/internal/entities"
	crud_services "berlin.allobank.com/tools/bonclay/internal/services/crud-services"
)

type CustomVariableService interface {
	// Legacy methods kept for backward compatibility
	GetCustomVariable(ctx context.Context, slug string) ([]repository2.CustomVariableEntity, error)
	CreateOrUpdateCustomVariable(ctx context.Context, id string, req CustomVariableReq) (repository2.CustomVariableEntity, error)

	// v2 methods
	ListVariables(ctx context.Context, collectionSlug string) ([]repository2.CustomVariableEntity, error)
	CreateVariable(ctx context.Context, collectionID, key, value string) (repository2.CustomVariableEntity, error)
	UpdateVariable(ctx context.Context, id, key, value string) (repository2.CustomVariableEntity, error)
	DeleteVariable(ctx context.Context, id string) error
}

type customVariableService struct {
	collectionService crud_services.BaseCRUDService[entities.CollectionEntity]
	repoContainer     repository2.RepoContainer
}

func NewCustomVariableService(
	collectionService crud_services.BaseCRUDService[entities.CollectionEntity],
	repoContainer repository2.RepoContainer,
) CustomVariableService {
	return &customVariableService{
		collectionService: collectionService,
		repoContainer:     repoContainer,
	}
}

func (s *customVariableService) CreateOrUpdateCustomVariable(ctx context.Context, id string, req CustomVariableReq) (repository2.CustomVariableEntity, error) {
	switch req.Method {
	case http.MethodGet, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodPost,
		http.MethodOptions, http.MethodTrace:
	default:
		return repository2.CustomVariableEntity{}, errs.BadRequest("invalid method")
	}

	resp, err := s.repoContainer.GetCustomVariableRepo().CreateOrUpdateCustomVariable(ctx, id, req.translate())
	if err != nil {
		return repository2.CustomVariableEntity{}, err
	}
	return resp, nil
}

func (s *customVariableService) GetCustomVariable(ctx context.Context, slug string) ([]repository2.CustomVariableEntity, error) {
	resp := []repository2.CustomVariableEntity{}

	// Get collection by slug using CRUD service
	collQuery := &pagination.ListQuery{ShowAll: true, Filters: map[string]interface{}{"slug = ?": slug}}
	collResult, err := s.collectionService.GetList(ctx, collQuery)
	if err != nil || len(collResult.List) == 0 {
		return []repository2.CustomVariableEntity{}, fmt.Errorf("collection not found")
	}
	collection := collResult.List[0]

	resp, err = s.repoContainer.GetCustomVariableRepo().GetListCustomVariableByCollectionId(ctx, collection.ID.String())
	if err != nil {
		return []repository2.CustomVariableEntity{}, err
	}

	return resp, nil
}

func (req CustomVariableReq) translate() repository2.CustomVariableEntity {
	return repository2.CustomVariableEntity{
		CollectionID: req.CollectionID,
		Key:          req.Key,
		Value:        req.Value,
	}
}

func (s *customVariableService) ListVariables(ctx context.Context, collectionSlug string) ([]repository2.CustomVariableEntity, error) {
	collQuery := &pagination.ListQuery{ShowAll: true, Filters: map[string]interface{}{"slug = ?": collectionSlug}}
	collResult, err := s.collectionService.GetList(ctx, collQuery)
	if err != nil || len(collResult.List) == 0 {
		return nil, fmt.Errorf("collection not found")
	}
	return s.repoContainer.GetCustomVariableRepo().GetListCustomVariableByCollectionId(ctx, collResult.List[0].ID.String())
}

func (s *customVariableService) CreateVariable(ctx context.Context, collectionID, key, value string) (repository2.CustomVariableEntity, error) {
	entity := repository2.CustomVariableEntity{
		CollectionID: collectionID,
		Key:          key,
		Value:        value,
	}
	return s.repoContainer.GetCustomVariableRepo().CreateOrUpdateCustomVariable(ctx, "", entity)
}

func (s *customVariableService) UpdateVariable(ctx context.Context, id, key, value string) (repository2.CustomVariableEntity, error) {
	existing, err := s.repoContainer.GetCustomVariableRepo().GetCustomVariableById(ctx, id)
	if err != nil {
		return repository2.CustomVariableEntity{}, err
	}
	existing.Key = key
	existing.Value = value
	return s.repoContainer.GetCustomVariableRepo().CreateOrUpdateCustomVariable(ctx, id, existing)
}

func (s *customVariableService) DeleteVariable(ctx context.Context, id string) error {
	return s.repoContainer.GetCustomVariableRepo().DeleteCustomVariableById(ctx, id)
}
