package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const customVariable = "custom_variables"

// CustomVariableEntity - domain entity for custom variables
type CustomVariableEntity struct {
	BaseEntityWithID
	CollectionID string `json:"collection_id" gorm:"column:collection_id"`
	Key          string `json:"key" gorm:"key"`
	Value        string `json:"value" gorm:"value"`
}

func (CustomVariableEntity) TableName() string {
	return customVariable
}

// CustomVariableRepoInterface - custom operations for custom variables
type CustomVariableRepoInterface interface {
	CreateOrUpdateCustomVariable(ctx context.Context, customVariableID string, entity CustomVariableEntity) (CustomVariableEntity, error)
	GetCustomVariableById(ctx context.Context, customVariableID string) (CustomVariableEntity, error)
	GetCustomVariableByKeyAndCollectionId(ctx context.Context, collectionId, key string) (CustomVariableEntity, error)
	GetListCustomVariableByCollectionId(ctx context.Context, collectionId string) ([]CustomVariableEntity, error)
}

type customVariableRepository struct {
	dbClient DBClient
}

func NewCustomVariableRepo(dbClient DBClient) CustomVariableRepoInterface {
	return &customVariableRepository{dbClient: dbClient}
}

func (r *customVariableRepository) CreateOrUpdateCustomVariable(ctx context.Context, customVariableID string, entity CustomVariableEntity) (CustomVariableEntity, error) {
	exec := r.dbClient.client(ctx).Table(customVariable)

	if customVariableID == "" {
		exec = exec.Create(&entity)
	} else {
		exec = exec.Where("id=?", customVariableID).Updates(&entity)
	}

	if exec.Error != nil {
		return CustomVariableEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	if exec.RowsAffected < 1 {
		return CustomVariableEntity{}, fmt.Errorf("cannot create data")
	}

	return entity, nil
}

func (r *customVariableRepository) GetCustomVariableById(ctx context.Context, customVariableID string) (CustomVariableEntity, error) {
	var entity CustomVariableEntity

	exec := r.dbClient.client(ctx).Table(customVariable).
		Where("id=?", customVariableID).
		First(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return CustomVariableEntity{}, fmt.Errorf("custom variable not found")
		}
		return CustomVariableEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (r *customVariableRepository) GetCustomVariableByKeyAndCollectionId(ctx context.Context, collectionId, key string) (CustomVariableEntity, error) {
	var entity CustomVariableEntity

	exec := r.dbClient.client(ctx).Table(customVariable).
		Where("collection_id=?", collectionId).
		Where("key=?", key).
		First(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return CustomVariableEntity{}, fmt.Errorf("custom variable not found")
		}
		return CustomVariableEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (r *customVariableRepository) GetListCustomVariableByCollectionId(ctx context.Context, collectionId string) ([]CustomVariableEntity, error) {
	var entity []CustomVariableEntity

	exec := r.dbClient.client(ctx).Model(&CustomVariableEntity{}).
		Where("collection_id=?", collectionId).
		Find(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return []CustomVariableEntity{}, fmt.Errorf("custom variable not found")
		}
		return []CustomVariableEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}
