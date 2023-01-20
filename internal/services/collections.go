package services

import (
	"context"
	"gitlab.linkaja.com/be/ditto/internal/repository"
)

type CollectionEntityReq struct {
	Name    string `json:"name" validate:"required"`
	Docs    string `json:"docs"`
	SquadID string `json:"squad_id" validate:"required,uuid"`
	Desc    string `json:"desc"`
}

func (req CollectionEntityReq) translate() repository.CollectionEntity {
	var ent repository.CollectionEntity
	if req.Name != "" {
		ent.Name = req.Name
	}
	if req.Docs != "" {
		ent.Docs = req.Name
	}
	if req.Desc != "" {
		ent.Desc = req.Desc
	}
	if req.SquadID != "" {
		ent.SquadID = req.SquadID
	}
	return ent
}

type CollectionServiceInterface interface {
	GetCollection(ctx context.Context, slug string) (repository.CollectionEntity, error)
	CreateOrUpdateCollection(ctx context.Context, squadID string, req CollectionEntityReq) (repository.CollectionEntity, error)
	DeleteCollection(ctx context.Context, squadID string) (repository.CollectionEntity, error)
}

func (cont serviceContainer) GetCollection(ctx context.Context, slug string) (repository.CollectionEntity, error) {
	resp, err := cont.repoContainer.GetCollectionBySlug(ctx, slug)
	if err != nil {
		return repository.CollectionEntity{}, err
	}

	return resp, nil
}

func (cont serviceContainer) CreateOrUpdateCollection(ctx context.Context, squadID string, req CollectionEntityReq) (repository.CollectionEntity, error) {
	resp, err := cont.repoContainer.CreateOrUpdateCollection(ctx, squadID, req.translate())
	if err != nil {
		return repository.CollectionEntity{}, err
	}
	return resp, nil
}

func (cont serviceContainer) DeleteCollection(ctx context.Context, squadID string) (repository.CollectionEntity, error) {
	_, err := cont.repoContainer.GetCollection(ctx, squadID)
	if err != nil {
		return repository.CollectionEntity{}, err
	}

	entity, err := cont.repoContainer.DeleteCollection(ctx, squadID)
	if err != nil {
		return repository.CollectionEntity{}, err
	}

	return entity, nil
}
