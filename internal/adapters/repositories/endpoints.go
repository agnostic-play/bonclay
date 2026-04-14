package repositories

import (
	"context"
	"errors"
	"fmt"

	"berlin.allobank.com/tools/bonclay/internal/entities"
	"gorm.io/gorm"
)

// EndpointRepoInterface - extends BaseCRUDRepo with custom methods
type EndpointRepoInterface interface {
	BaseCRUDRepo[entities.EndpointEntity]
	GetEndpointMock(ctx context.Context, collectionID, method, path string) (entities.EndpointEntity, error)
	GetListEndpointByCollectionID(ctx context.Context, collectionID, method string) ([]entities.EndpointEntity, error)
}

type endpointRepository struct {
	*baseRepository[entities.EndpointEntity]
}

func NewEndpointRepo(dbClient DBClient) EndpointRepoInterface {
	ent := &entities.EndpointEntity{}
	return &endpointRepository{
		baseRepository: &baseRepository[entities.EndpointEntity]{
			ent:      ent,
			dbClient: dbClient,
		},
	}
}

func (r *endpointRepository) GetEndpointMock(ctx context.Context, collectionID, method, path string) (entities.EndpointEntity, error) {
	var entity entities.EndpointEntity

	exec := r.dbClient.client(ctx).Table("endpoints").
		Where("collection_id=?", collectionID).
		Where("method=? AND path=?", method, path).
		First(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return entities.EndpointEntity{}, fmt.Errorf("endpoint not found")
		}
		return entities.EndpointEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}

func (r *endpointRepository) GetListEndpointByCollectionID(ctx context.Context, collectionID, method string) ([]entities.EndpointEntity, error) {
	var entity []entities.EndpointEntity

	exec := r.dbClient.client(ctx).Table("endpoints").
		Where("collection_id=?", collectionID).
		Where("method=?", method).
		Find(&entity)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return []entities.EndpointEntity{}, fmt.Errorf("endpoint not found")
		}

		return []entities.EndpointEntity{}, fmt.Errorf("query error: %s", exec.Error)
	}

	return entity, nil
}
