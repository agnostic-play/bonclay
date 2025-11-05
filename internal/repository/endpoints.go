package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const endpoints = "endpoints"

type EndpointEntity struct {
	BaseEntityWithID

	Path           string `json:"path" gorm:"column:path"`
	Method         string `json:"method" gorm:"column:method"`
	CollectionID   string `json:"collection_id" gorm:"column:collection_id"`
	Category       string `json:"category" gorm:"column:category"`
	ActiveScenario string `json:"active_scenario" gorm:"column:active_scenario"`
	Desc           string `json:"desc" gorm:"column:desc"`

	EnableResponseIntercept bool `json:"enable_response_intercept" gorm:"column:enable_response_intercept"`
}

func (EndpointEntity) TableName() string {
	return "endpoints"
}

type EndpointRepoInterface interface {
	GetEndpoint(ctx context.Context, id string) (EndpointEntity, error)
	GetEndpointMock(ctx context.Context, collectionID, method, path string) (EndpointEntity, error)
	CreateOrUpdateEndpoint(ctx context.Context, squadID string, entity EndpointEntity) (EndpointEntity, error)
	DeleteEndpoint(ctx context.Context, squadID string) error
}

func (cont repoContainerGorm) CreateOrUpdateEndpoint(ctx context.Context, endpointID string, entity EndpointEntity) (EndpointEntity, error) {
	exec := cont.write(ctx).Table(endpoints)

	if endpointID == "" {
		exec = exec.Omit("active_scenario").Create(&entity)
	} else {
		exec = exec.Omit("collection_id", "active_scenario").Where("id=?", endpointID).Updates(&entity)
	}

	if exec.Error != nil {
		return EndpointEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	if exec.RowsAffected < 1 {
		return EndpointEntity{}, fmt.Errorf("cannot create data")
	}

	return entity, nil
}

func (cont repoContainerGorm) GetEndpointMock(ctx context.Context, collectionID, method, path string) (EndpointEntity, error) {

	var entity EndpointEntity

	exec := cont.db.WithContext(ctx).Table(endpoints).
		Where("collection_id=?", collectionID).
		Where("method=? AND path=?", method, path).
		First(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return EndpointEntity{}, fmt.Errorf("collections not found")
		}
		return EndpointEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (cont repoContainerGorm) GetEndpoint(ctx context.Context, slug string) (EndpointEntity, error) {

	var entity EndpointEntity

	exec := cont.db.WithContext(ctx).Table(endpoints).Where("id=?", slug).First(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return EndpointEntity{}, fmt.Errorf("collections not found")
		}
		return EndpointEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (cont repoContainerGorm) DeleteEndpoint(ctx context.Context, squadID string) error {

	var entity EndpointEntity

	exec := cont.write(ctx).Table(endpoints).Where("id=?", squadID).Delete(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("collections not found")
		}
		return fmt.Errorf("query error: %s", exec.Error)
	}

	return nil
}
