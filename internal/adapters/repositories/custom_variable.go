package repositories

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
	CreateOrUpdateMultipleCustomVariable(ctx context.Context, collectionId string, entities map[string]string) (affectedRow int, err error)
	DeleteCustomVariableById(ctx context.Context, id string) error
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

// CreateOrUpdateMultipleCustomVariable upserts a batch of custom variables for a
// collection within a single transaction. Existing keys are updated in place;
// unknown keys are inserted. It is used to persist variables produced by an
// endpoint script after a mock request runs.
func (r *customVariableRepository) CreateOrUpdateMultipleCustomVariable(ctx context.Context, collectionId string, entities map[string]string) (affectedRow int, err error) {
	if len(entities) == 0 {
		return 0, nil
	}

	tx := r.dbClient.client(ctx).Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("begin tx error: %w", tx.Error)
	}

	for key, value := range entities {
		var existing CustomVariableEntity

		findErr := tx.Table(customVariable).
			Where("collection_id = ? AND `key` = ?", collectionId, key).
			First(&existing).Error

		switch {
		case errors.Is(findErr, gorm.ErrRecordNotFound):
			newEntity := CustomVariableEntity{
				CollectionID: collectionId,
				Key:          key,
				Value:        value,
			}
			if createErr := tx.Table(customVariable).Create(&newEntity).Error; createErr != nil {
				tx.Rollback()
				return affectedRow, fmt.Errorf("create error for key=%s: %w", key, createErr)
			}
			affectedRow++
		case findErr == nil:
			exec := tx.Table(customVariable).
				Where("id = ?", existing.ID).
				Update("value", value)
			if exec.Error != nil {
				tx.Rollback()
				return affectedRow, fmt.Errorf("update error for key=%s: %w", key, exec.Error)
			}
			if exec.RowsAffected > 0 {
				affectedRow++
			}
		default:
			tx.Rollback()
			return affectedRow, fmt.Errorf("query error for key=%s: %w", key, findErr)
		}
	}

	if commitErr := tx.Commit().Error; commitErr != nil {
		return affectedRow, fmt.Errorf("commit error: %w", commitErr)
	}

	return affectedRow, nil
}

func (r *customVariableRepository) DeleteCustomVariableById(ctx context.Context, id string) error {
	exec := r.dbClient.client(ctx).Table(customVariable).
		Where("id = ?", id).
		Delete(&CustomVariableEntity{})
	if exec.Error != nil {
		return fmt.Errorf("query error: %s", exec.Error)
	}
	if exec.RowsAffected < 1 {
		return fmt.Errorf("custom variable not found")
	}
	return nil
}
