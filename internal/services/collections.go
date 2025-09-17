package services

import (
	"context"
	"strings"

	"github.com/agnostic-play/ditoo/internal/repository"
)

type CollectionEntityReq struct {
	Name    string `json:"name" validate:"required"`
	Docs    string `json:"docs"`
	SquadID string `json:"squad_id" validate:"required,uuid"`
	Desc    string `json:"desc"`

	ForwardProxyURL string `json:"forward_proxy_url"`
	IsProxyEnable   bool   `json:"is_proxy_enable" `
}

func (req CollectionEntityReq) translate() repository.CollectionEntity {
	var ent repository.CollectionEntity

	if val := strings.TrimSpace(req.Name); val != "" {
		ent.Name = val
	}

	if val := strings.TrimSpace(req.Docs); val != "" {
		ent.Docs = val
	}

	if val := strings.TrimSpace(req.Desc); val != "" {
		ent.Desc = val
	}

	if val := strings.TrimSpace(req.ForwardProxyURL); val != "" {
		ent.ForwardProxyURL = val
	}

	ent.IsProxyEnable = req.IsProxyEnable
	ent.SquadID = req.SquadID
	return ent
}

type CollectionServiceInterface interface {
	GetCollection(ctx context.Context, slug string) (repository.CollectionEntity, error)
	CreateOrUpdateCollection(ctx context.Context, id string, req CollectionEntityReq) (repository.CollectionEntity, error)
	DeleteCollection(ctx context.Context, id string) (repository.CollectionEntity, error)
}

func (cont serviceContainer) GetCollection(ctx context.Context, slug string) (repository.CollectionEntity, error) {
	resp, err := cont.repoContainer.GetCollectionBySlug(ctx, slug)
	if err != nil {
		return repository.CollectionEntity{}, err
	}

	return resp, nil
}

func (cont serviceContainer) CreateOrUpdateCollection(ctx context.Context, id string, req CollectionEntityReq) (repository.CollectionEntity, error) {
	resp, err := cont.repoContainer.CreateOrUpdateCollection(ctx, id, req.translate())
	if err != nil {
		return repository.CollectionEntity{}, err
	}
	return resp, nil
}

func (cont serviceContainer) DeleteCollection(ctx context.Context, id string) (repository.CollectionEntity, error) {
	_, err := cont.repoContainer.GetCollection(ctx, id)
	if err != nil {
		return repository.CollectionEntity{}, err
	}

	entity, err := cont.repoContainer.DeleteCollection(ctx, id)
	if err != nil {
		return repository.CollectionEntity{}, err
	}

	return entity, nil
}
