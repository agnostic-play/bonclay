package crud_services

import (
	"context"
	"errors"

	repository2 `github.com/agnostic-play/ditoo/internal/adapters/repositories`
	"github.com/agnostic-play/ditoo/internal/common/errs"
	"github.com/agnostic-play/ditoo/internal/common/pagination"
	"github.com/agnostic-play/ditoo/internal/entities"
	"gorm.io/gorm"
)

type BaseCRUDService[T any] interface {
	GetList(ctx context.Context, query *pagination.ListQuery) (*pagination.ListResult[T], error)
	Create(ctx context.Context, entity T) (T, error)

	Get(ctx context.Context, id string) (T, error)
	Patch(ctx context.Context, id string, entity T) (T, error)
	Delete(ctx context.Context, id string) error

	GetEntity() entities.BaseEntityInterface[T]
}

type baseCRUDServices[T any] struct {
	ent  entities.BaseEntityInterface[T]
	repo repository2.BaseCRUDRepo[T]
}

func NewBaseCRUDServices[T any](entity entities.BaseEntityInterface[T], dbClient repository2.DBClient) BaseCRUDService[T] {
	return &baseCRUDServices[T]{
		ent:  entity,
		repo: repository2.NewBaseCRUDRepo(entity, dbClient),
	}
}

func (b baseCRUDServices[T]) GetList(ctx context.Context, query *pagination.ListQuery) (*pagination.ListResult[T], error) {
	result, err := b.repo.GetList(ctx, query)
	if err != nil {
		err = errs.DatabaseErr(err)
		return nil, err
	}

	return result, nil
}

func (b baseCRUDServices[T]) Create(ctx context.Context, entity T) (T, error) {
	var result T

	result, err := b.repo.Create(ctx, entity)
	if err != nil {
		err = errs.DatabaseErr(err)
		return result, err
	}

	return result, nil
}

func (b baseCRUDServices[T]) Get(ctx context.Context, id string) (T, error) {
	var result T

	result, err := b.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errs.NotFoundErr(b.ent.GetTitle())
			return result, nil
		}

		err = errs.DatabaseErr(err)
		return result, err
	}

	return result, nil
}

func (b baseCRUDServices[T]) Patch(ctx context.Context, id string, entity T) (T, error) {
	var result T

	result, err := b.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errs.NotFoundErr(b.ent.GetTitle())
			return result, nil
		}

		err = errs.DatabaseErr(err)
		return result, err
	}

	result, err = b.repo.Patch(ctx, id, entity)
	if err != nil {
		err = errs.DatabaseErr(err)
		return result, err
	}

	return result, nil
}

func (b baseCRUDServices[T]) Delete(ctx context.Context, id string) error {
	_, err := b.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errs.NotFoundErr(b.ent.GetTitle())
			return nil
		}

		err = errs.DatabaseErr(err)
		return err
	}

	if err = b.repo.Delete(ctx, id); err != nil {
		err = errs.DatabaseErr(err)
		return err
	}

	return nil
}

func (b baseCRUDServices[T]) GetEntity() entities.BaseEntityInterface[T] {
	return b.ent
}
