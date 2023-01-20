package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type SquadEntity struct {
	BaseEntityWithID
	Name        string             `json:"name"`
	Slug        string             `json:"slug"`
	Password    string             `json:"password"`
	Desc        string             `json:"desc"`
	Collections []CollectionEntity `json:"collections,omitempty" gorm:"-"`
}

type SquadServiceInterface interface {
	GetSquad(ctx context.Context, slug string) (SquadEntity, error)
	GetSquadBySlug(ctx context.Context, slug string) (SquadEntity, error)
	GetListSquad(ctx context.Context, filterName string) ([]SquadEntity, error)
	CreateOrUpdateSquad(ctx context.Context, squadID string, entity SquadEntity) (SquadEntity, error)
	DeleteSquad(ctx context.Context, squadID string) (SquadEntity, error)
}

func (cont repoContainerGorm) CreateOrUpdateSquad(ctx context.Context, squadID string, entity SquadEntity) (SquadEntity, error) {

	exec := cont.write(ctx).Table("squads")
	if squadID == "" {
		entity.Slug = cont.generateSlug("squads", entity.Name)
		exec = exec.Create(&entity)
	} else {
		exec = exec.Omit("slug").Where("id=?", squadID).Updates(&entity)
	}

	if exec.Error != nil {
		return SquadEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	if exec.RowsAffected < 1 {
		return SquadEntity{}, fmt.Errorf("cannot create data")
	}

	return entity, nil
}

func (cont repoContainerGorm) GetListSquad(ctx context.Context, filterName string) ([]SquadEntity, error) {

	var entity []SquadEntity

	exec := cont.db.WithContext(ctx).Table("squads")
	if filterName != "" {
		exec = exec.Where("name LIKE ?", "%"+filterName+"%")
	}
	exec = exec.Omit("password").Omit("collections").Find(&entity)

	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return []SquadEntity{}, fmt.Errorf("squad not found")
		}
		return []SquadEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (cont repoContainerGorm) GetSquad(ctx context.Context, slug string) (SquadEntity, error) {

	var entity SquadEntity

	exec := cont.db.WithContext(ctx).Table("squads").Omit("passwords").Where("id=?", slug).First(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return SquadEntity{}, fmt.Errorf("squad not found")
		}
		return SquadEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (cont repoContainerGorm) GetSquadBySlug(ctx context.Context, slug string) (SquadEntity, error) {

	var entity SquadEntity

	exec := cont.db.WithContext(ctx).Table("squads").Omit("passwords").Where("slug=?", slug).First(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return SquadEntity{}, fmt.Errorf("squad not found")
		}
		return SquadEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (cont repoContainerGorm) DeleteSquad(ctx context.Context, squadID string) (SquadEntity, error) {

	var entity SquadEntity

	exec := cont.write(ctx).Table("squads").Where("id=?", squadID).Delete(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return SquadEntity{}, fmt.Errorf("squad not found")
		}
		return SquadEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}
