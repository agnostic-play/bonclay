package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const customVariable = "custom_variables"

type CustomVariableEntity struct {
	BaseEntityWithID
	CollectionID string `json:"collection_id" gorm:"column:collection_id"`
	Key          string `json:"key" gorm:"key"`
	Value        string `json:"value" gorm:"value"`
}

func (CustomVariableEntity) TableName() string {
	return "custom_variables"
}

type CustomVariableRepoInterface interface {
	CreateOrUpdateCustomVariable(ctx context.Context, customVariableID string, entity CustomVariableEntity) (CustomVariableEntity, error)
	GetCustomVariableById(ctx context.Context, customVariableID string) (CustomVariableEntity, error)
	GetCustomVariableByKeyAndCollectionId(ctx context.Context, collectionId, key string) (CustomVariableEntity, error)
	GetListCustomVariableByCollectionId(ctx context.Context, collectionId string) ([]CustomVariableEntity, error)
	CreateOrUpdateMultipleCustomVariable(ctx context.Context, collectionId string, entities map[string]string) (affectedRow int, err error)
}

func (cont repoContainerGorm) CreateOrUpdateCustomVariable(ctx context.Context, customVariableID string, entity CustomVariableEntity) (CustomVariableEntity, error) {
	exec := cont.write(ctx).Table(customVariable)

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

func (cont repoContainerGorm) GetCustomVariableById(ctx context.Context, customVariableID string) (CustomVariableEntity, error) {
	var entity CustomVariableEntity

	exec := cont.db.WithContext(ctx).Table(customVariable).
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

func (cont repoContainerGorm) GetCustomVariableByKeyAndCollectionId(ctx context.Context, collectionId, key string) (CustomVariableEntity, error) {
	var entity CustomVariableEntity

	exec := cont.db.WithContext(ctx).Table(customVariable).
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

func (cont repoContainerGorm) GetListCustomVariableByCollectionId(ctx context.Context, collectionId string) ([]CustomVariableEntity, error) {
	var entity []CustomVariableEntity

	exec := cont.db.Model(&CustomVariableEntity{}).
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

func (cont repoContainerGorm) CreateOrUpdateMultipleCustomVariable(ctx context.Context, collectionId string, entities map[string]string) (affectedRow int, err error) {
	if len(entities) == 0 {
		return 0, nil
	}

	tx := cont.write(ctx).Begin()

	for key, value := range entities {
		var existing CustomVariableEntity

		// check if variable already exists in this collection
		err := tx.Table(customVariable).
			Where("collection_id = ? AND `key` = ?", collectionId, key).
			First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create new
			newEntity := CustomVariableEntity{
				CollectionID: collectionId,
				Key:          key,
				Value:        value,
			}
			if err := tx.Table(customVariable).Create(&newEntity).Error; err != nil {
				tx.Rollback()
				return affectedRow, fmt.Errorf("create error for key=%s: %w", key, err)
			}
			affectedRow++
		} else if err == nil {
			// update existing
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
		} else {
			// unexpected error during lookup
			tx.Rollback()
			return affectedRow, fmt.Errorf("query error for key=%s: %w", key, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return affectedRow, fmt.Errorf("commit error: %w", err)
	}

	return affectedRow, nil
}
