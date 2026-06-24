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
	UpdateScript(ctx context.Context, endpointID, script string) error
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

// UpdateScript sets (or clears) the JavaScript snippet attached to an endpoint.
// A direct column update is used so an empty string can be persisted, unlike the
// generic patch which skips zero-valued fields.
func (r *endpointRepository) UpdateScript(ctx context.Context, endpointID, script string) error {
	exec := r.dbClient.client(ctx).Table("endpoints").
		Where("id = ?", endpointID).
		Update("script", script)
	if exec.Error != nil {
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("endpoint not found")
		}
		return fmt.Errorf("query error: %s", exec.Error)
	}
	if exec.RowsAffected < 1 {
		return fmt.Errorf("endpoint not found")
	}

	return nil
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
