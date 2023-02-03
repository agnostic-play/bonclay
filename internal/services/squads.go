package services

import (
	"context"
	"gitlab.linkaja.com/be/ditto/internal/repository"
	"strings"
)

type SquadEntityReq struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password,omitempty" validate:""`
	Desc     string `json:"desc"`
}

func (req SquadEntityReq) translate() repository.SquadEntity {
	var ent repository.SquadEntity

	if val := strings.TrimSpace(req.Name); val != "" {
		ent.Name = val
	}

	if val := strings.TrimSpace(req.Password); val != "" {
		ent.Password = val
	}

	if val := strings.TrimSpace(req.Desc); val != "" {
		ent.Desc = val
	}

	return ent
}

type SquadServiceInterface interface {
	GetSquad(ctx context.Context, slug string) (repository.SquadEntity, error)
	GetListSquad(ctx context.Context, name string) ([]repository.SquadEntity, error)
	CreateOrUpdateSquad(ctx context.Context, id string, req SquadEntityReq) (repository.SquadEntity, error)
	DeleteSquad(ctx context.Context, id string) (repository.SquadEntity, error)
}

func (cont serviceContainer) GetListSquad(ctx context.Context, name string) ([]repository.SquadEntity, error) {
	resp, err := cont.repoContainer.GetListSquad(ctx, name)
	if err != nil {
		return []repository.SquadEntity{}, err
	}

	return resp, nil
}

func (cont serviceContainer) GetSquad(ctx context.Context, slug string) (repository.SquadEntity, error) {
	resp, err := cont.repoContainer.GetSquadBySlug(ctx, slug)
	if err != nil {
		return repository.SquadEntity{}, err
	}

	if collections, err := cont.repoContainer.GetListCollection(ctx, resp.ID.String()); err == nil {
		resp.Collections = collections
	}

	return resp, nil
}

func (cont serviceContainer) CreateOrUpdateSquad(ctx context.Context, id string, req SquadEntityReq) (repository.SquadEntity, error) {
	resp, err := cont.repoContainer.CreateOrUpdateSquad(ctx, id, req.translate())
	if err != nil {
		return repository.SquadEntity{}, err
	}
	resp.Password = ""
	return resp, nil
}

func (cont serviceContainer) DeleteSquad(ctx context.Context, id string) (repository.SquadEntity, error) {
	_, err := cont.repoContainer.GetSquad(ctx, id)
	if err != nil {
		return repository.SquadEntity{}, err
	}

	entity, err := cont.repoContainer.DeleteSquad(ctx, id)
	if err != nil {
		return repository.SquadEntity{}, err
	}

	return entity, nil
}
