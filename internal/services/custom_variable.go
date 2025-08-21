package services

import (
	"context"

	"github.com/agnostic-play/ditoo/internal/repository"
)

type CustomVariableReq struct {
	CollectionID string `json:"collection_id" validate:"required,uuid"`
	Key          string `json:"key" validate:"required"`
	Value        string `json:"value" validate:"required"`
}

type CustomVariableCreateReq struct {
	CustomVariableReq
	CollectionID string `json:"collection_id" validate:"required,uuid"`
}

func (req CustomVariableReq) translate() repository.CustomVariableEntity {
	return repository.CustomVariableEntity{
		CollectionID: req.CollectionID,
		Key:          req.Key,
		Value:        req.Value,
	}
}

type CustomVariableServiceInterface interface {
	GetCustomVariable(ctx context.Context, slug string) ([]repository.CustomVariableEntity, error)
}

func (cont serviceContainer) GetCustomVariable(ctx context.Context, slug string) ([]repository.CustomVariableEntity, error) {

	resp := []repository.CustomVariableEntity{}

	collection, err := cont.repoContainer.GetCollectionBySlug(ctx, slug)
	if err != nil {
		return []repository.CustomVariableEntity{}, err
	}

	resp, err = cont.repoContainer.GetListCustomVariableByCollectionId(ctx, collection.ID.String())
	if err != nil {
		return []repository.CustomVariableEntity{}, err
	}

	return resp, nil
}
