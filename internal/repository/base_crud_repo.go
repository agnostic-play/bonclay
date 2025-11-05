package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/agnostic-play/ditoo/internal/common/pagination"
	"github.com/agnostic-play/ditoo/internal/entities"
	"gorm.io/gorm"
)

type BaseCRUDRepo[T any] interface {
	GetList(ctx context.Context, query *pagination.ListQuery) (*pagination.ListResult[T], error)
	Create(ctx context.Context, entity T) (T, error)

	Get(ctx context.Context, id string) (T, error)
	Patch(ctx context.Context, id string, entity T) (T, error)
	Delete(ctx context.Context, id string) error
}

type baseRepository[T any] struct {
	ent      entities.BaseEntityInterface[T]
	dbClient DBClient
}

func NewBaseCRUDRepo[T any](entity entities.BaseEntityInterface[T], dbClient DBClient) BaseCRUDRepo[T] {
	return &baseRepository[T]{
		ent:      entity,
		dbClient: dbClient,
	}
}

func (r *baseRepository[T]) GetList(ctx context.Context, paginationQuery *pagination.ListQuery) (*pagination.ListResult[T], error) {
	var (
		totalItems int64
		listItems  []T
		result     = pagination.NewListResult[T](paginationQuery)
	)

	log.Println("xxxxx")
	rawExec := r.dbClient.client(ctx).
		Model(r.ent.GetEntity()).
		Table(r.ent.GetTableName()).
		Scopes(SearchScope(paginationQuery, r.ent.GetFieldForKeywords()...))

	getCount := rawExec.Count(&totalItems)
	if err := getCount.Error; err != nil {
		return &pagination.ListResult[T]{}, err
	}

	if totalItems < 1 {
		return result, nil
	}

	getItems := rawExec.Scopes(Sort(paginationQuery))

	if !paginationQuery.ShowAll {
		getItems = getItems.Scopes(Paginate(paginationQuery))
	}

	getItems = getItems.Find(&listItems)
	if err := getItems.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, nil
		}

		return &pagination.ListResult[T]{}, err
	}

	result.SetData(int(totalItems), listItems)

	return result, nil
}

func (r *baseRepository[T]) Create(ctx context.Context, entity T) (T, error) {
	var result T

	exec := r.dbClient.client(ctx).
		Model(r.ent.GetEntity()).
		Table(r.ent.GetTableName()).
		Create(&entity)

	if exec.Error != nil {
		return result, exec.Error
	}

	if exec.RowsAffected < 1 {
		return result, fmt.Errorf("%s | failed to insert entity", r.ent.GetTitle())
	}

	return result, nil
}

func (r *baseRepository[T]) Get(ctx context.Context, id string) (T, error) {
	var result T

	exec := r.dbClient.client(ctx).
		Model(r.ent.GetEntity()).
		Table(r.ent.GetTableName()).
		Where(fmt.Sprintf("%s = ?", r.ent.GetFieldIdentifier()), id).
		First(&result)
	if exec.Error != nil {
		return result, exec.Error
	}

	return result, nil
}

func (r *baseRepository[T]) Patch(ctx context.Context, id string, entity T) (T, error) {
	var result T

	// Perform partial update (only the fields in 'updates' will be modified)
	exec := r.dbClient.client(ctx).
		Table(r.ent.GetTableName()).
		Model(r.ent.GetEntity()).
		Where("id = ?", id).
		Omit(r.ent.GetExcludeFieldForUpdate()...).
		Updates(r.structToMap(entity))
	if exec.Error != nil {
		return result, exec.Error
	}

	if exec.RowsAffected < 1 {
		return result, fmt.Errorf("%s | failed to update entity", r.ent.GetTitle())
	}

	return result, nil
}

// Delete deletes an entity by its ID
func (r *baseRepository[T]) Delete(ctx context.Context, id string) error {
	exec := r.dbClient.client(ctx).
		Model(new(T)).Table(r.ent.GetTableName()).
		Where("id = ?", id).Delete(new(T))
	if exec.Error != nil {
		return exec.Error
	}

	if exec.RowsAffected < 1 {
		return fmt.Errorf("%s | failed to delete entity", r.ent.GetTitle())
	}

	return nil
}

func (r *baseRepository[T]) structToMap(entity T) map[string]interface{} {
	updates := make(map[string]interface{})
	val := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)

	// Ensure the input is a struct and not a pointer
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Iterate over the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)

		// Only include non-zero fields
		if !fieldValue.IsZero() {
			jsonTag := fieldType.Tag.Get("json")
			fieldName := strings.Split(jsonTag, ",")[0]
			if fieldName == "" {
				fieldName = fieldType.Name
			}

			updates[fieldName] = fieldValue.Interface()
		}
	}

	return updates
}
