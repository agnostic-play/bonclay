package api_mock_services

import (
	"context"
	"fmt"
	"net/http"

	repository2 `github.com/agnostic-play/ditoo/internal/adapters/repositories`
	"github.com/agnostic-play/ditoo/internal/common/errs"
	"github.com/agnostic-play/ditoo/internal/common/pagination"
	"github.com/agnostic-play/ditoo/internal/entities"
	crud_services "github.com/agnostic-play/ditoo/internal/services/crud-services"
)

type CustomVariableService interface {
	GetCustomVariable(ctx context.Context, slug string) ([]repository2.CustomVariableEntity, error)
	CreateOrUpdateCustomVariable(ctx context.Context, id string, req CustomVariableReq) (repository2.CustomVariableEntity, error)
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
