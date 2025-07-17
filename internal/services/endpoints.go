package services

import (
	"context"
	"net/http"
	"strings"

	"github.com/agnostic-play/ditoo/internal/errs"
	"github.com/agnostic-play/ditoo/internal/repository"
)

type EndpointEntityReq struct {
	Path         string `json:"path" validate:"required"`
	Method       string `json:"method" validate:"required,alpha"`
	CollectionID string `json:"collection_id" validate:"required,uuid"`
	Category     string `json:"category" validate:"required,customText"`
	Desc         string `json:"desc" validate:"required"`
}

type EndpointCreateReq struct {
	EndpointEntityReq
	CollectionID string `json:"collection_id" validate:"required,uuid"`
}

func (req EndpointEntityReq) translate() repository.EndpointEntity {
	var ent repository.EndpointEntity

	if val := strings.TrimSpace(req.Method); val != "" {
		ent.Method = strings.ToUpper(val)
	}

	if val := strings.TrimSpace(req.Path); val != "" {
		ent.Path = strings.ToLower(val)
	}

	if val := strings.TrimSpace(req.Desc); val != "" {
		ent.Desc = strings.ToLower(val)
	}

	if val := strings.TrimSpace(req.CollectionID); val != "" {
		ent.CollectionID = val
	}

	if val := strings.TrimSpace(req.Category); val != "" {
		ent.Category = strings.ToLower(val)
	}

	return ent
}

type EndpointServiceInterface interface {
	CreateOrUpdateEndpoint(ctx context.Context, id string, req EndpointEntityReq) (repository.EndpointEntity, error)
	DeleteEndpoint(ctx context.Context, id string) error
}

func (cont serviceContainer) CreateOrUpdateEndpoint(ctx context.Context, id string, req EndpointEntityReq) (repository.EndpointEntity, error) {
	switch req.Method {
	case http.MethodGet, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodPost,
		http.MethodOptions, http.MethodTrace:
	default:
		return repository.EndpointEntity{}, errs.BadRequest("invalid method")
	}

	resp, err := cont.repoContainer.CreateOrUpdateEndpoint(ctx, id, req.translate())
	if err != nil {
		return repository.EndpointEntity{}, err
	}
	return resp, nil
}

func (cont serviceContainer) DeleteEndpoint(ctx context.Context, id string) error {
	_, err := cont.repoContainer.GetEndpoint(ctx, id)
	if err != nil {
		return err
	}

	err = cont.repoContainer.DeleteEndpoint(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
