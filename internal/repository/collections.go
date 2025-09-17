package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const collections = "collections"

type CollectionEntity struct {
	BaseEntityWithID
	Name    string `json:"name"`
	Docs    string `json:"docs"`
	Slug    string `json:"slug"`
	SquadID string `json:"squad_id"`
	Desc    string `json:"desc"`

	ForwardProxyURL string `json:"forward_proxy_url" gorm:"column:forward_proxy_url"`
	IsProxyEnable   bool   `json:"is_proxy_enable" gorm:"column:is_proxy_enable"`
}

type CollectionRepoInterface interface {
	GetCollection(ctx context.Context, id string) (CollectionEntity, error)
	GetCollectionBySlug(ctx context.Context, slug string) (CollectionEntity, error)
	GetListCollection(ctx context.Context, squadID string) ([]CollectionEntity, error)
	CreateOrUpdateCollection(ctx context.Context, squadID string, entity CollectionEntity) (CollectionEntity, error)
	DeleteCollection(ctx context.Context, squadID string) (CollectionEntity, error)
}

func (cont repoContainerGorm) CreateOrUpdateCollection(ctx context.Context, collectionID string, entity CollectionEntity) (CollectionEntity, error) {
	exec := cont.write(ctx).Table(collections)

	if collectionID == "" {
		entity.Slug = cont.generateSlug(collections, entity.Name)
		exec = exec.Create(&entity)
	} else {
		exec = exec.Omit("slug", "squad_id").Where("id=?", collectionID).Updates(&entity)
	}

	if exec.Error != nil {
		return CollectionEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	if exec.RowsAffected < 1 {
		return CollectionEntity{}, fmt.Errorf("cannot create data")
	}

	return entity, nil
}

func (cont repoContainerGorm) GetListCollection(ctx context.Context, squadID string) ([]CollectionEntity, error) {

	var entity []CollectionEntity

	exec := cont.db.WithContext(ctx).Table(collections).Where("squad_id = ?", squadID).Find(&entity)

	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return []CollectionEntity{}, nil
		}
		return []CollectionEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (cont repoContainerGorm) GetCollection(ctx context.Context, slug string) (CollectionEntity, error) {

	var entity CollectionEntity

	exec := cont.db.WithContext(ctx).Table(collections).Where("id=?", slug).First(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return CollectionEntity{}, fmt.Errorf("collections not found")
		}
		return CollectionEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (cont repoContainerGorm) GetCollectionBySlug(ctx context.Context, slug string) (CollectionEntity, error) {

	var entity CollectionEntity

	exec := cont.db.WithContext(ctx).Table(collections).Where("slug=?", slug).First(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return CollectionEntity{}, fmt.Errorf("collections not found")
		}
		return CollectionEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (cont repoContainerGorm) DeleteCollection(ctx context.Context, squadID string) (CollectionEntity, error) {

	var entity CollectionEntity

	exec := cont.write(ctx).Table(collections).Where("id=?", squadID).Delete(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return CollectionEntity{}, fmt.Errorf("collections not found")
		}
		return CollectionEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}
