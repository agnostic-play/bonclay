package api_mock_services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/agnostic-play/ditoo/internal/common/errs"
	"github.com/agnostic-play/ditoo/internal/common/pagination"
	"github.com/agnostic-play/ditoo/internal/entities"
	"github.com/agnostic-play/ditoo/internal/repository"
	crud_services "github.com/agnostic-play/ditoo/internal/services/crud-services"
)

type CustomVariableService interface {
	GetCustomVariable(ctx context.Context, slug string) ([]repository.CustomVariableEntity, error)
	CreateOrUpdateCustomVariable(ctx context.Context, id string, req CustomVariableReq) (repository.CustomVariableEntity, error)
}

type customVariableService struct {
	collectionService crud_services.BaseCRUDService[entities.CollectionEntity]
	repoContainer     repository.RepoContainer
}

func NewCustomVariableService(
	collectionService crud_services.BaseCRUDService[entities.CollectionEntity],
	repoContainer repository.RepoContainer,
) CustomVariableService {
	return &customVariableService{
		collectionService: collectionService,
		repoContainer:     repoContainer,
	}
}

func (s *customVariableService) CreateOrUpdateCustomVariable(ctx context.Context, id string, req CustomVariableReq) (repository.CustomVariableEntity, error) {
	switch req.Method {
	case http.MethodGet, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodPost,
		http.MethodOptions, http.MethodTrace:
	default:
		return repository.CustomVariableEntity{}, errs.BadRequest("invalid method")
	}

	resp, err := s.repoContainer.GetCustomVariableRepo().CreateOrUpdateCustomVariable(ctx, id, req.translate())
	if err != nil {
		return repository.CustomVariableEntity{}, err
	}
	return resp, nil
}

func (s *customVariableService) GetCustomVariable(ctx context.Context, slug string) ([]repository.CustomVariableEntity, error) {
	resp := []repository.CustomVariableEntity{}

	// Get collection by slug using CRUD service
	collQuery := &pagination.ListQuery{ShowAll: true, Filters: map[string]interface{}{"slug = ?": slug}}
	collResult, err := s.collectionService.GetList(ctx, collQuery)
	if err != nil || len(collResult.List) == 0 {
		return []repository.CustomVariableEntity{}, fmt.Errorf("collection not found")
	}
	collection := collResult.List[0]

	resp, err = s.repoContainer.GetCustomVariableRepo().GetListCustomVariableByCollectionId(ctx, collection.ID.String())
	if err != nil {
		return []repository.CustomVariableEntity{}, err
	}

	return resp, nil
}

func (req CustomVariableReq) translate() repository.CustomVariableEntity {
	return repository.CustomVariableEntity{
		CollectionID: req.CollectionID,
		Key:          req.Key,
		Value:        req.Value,
	}
}
