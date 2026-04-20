package handler

import (
	"context"
	"fmt"
	"net/http"

	"berlin.allobank.com/tools/bonclay/internal/entities"
	"berlin.allobank.com/tools/bonclay/internal/services"
	crud_services "berlin.allobank.com/tools/bonclay/internal/services/crud-services"
	"github.com/labstack/echo/v4"
)

// CollectionHandlers extends BaseCRUDHandlers with collection-specific endpoints
type CollectionHandlers interface {
	BaseCRUDHandlers[entities.CollectionEntity]
}

type collectionHandlers struct {
	baseCRUDHandlers BaseCRUDHandlers[entities.CollectionEntity]
	endpointServices crud_services.BaseCRUDService[entities.EndpointEntity]
	serviceContainer services.ServiceContainer
}

// NewCollectionHandlers creates a new CollectionRoutes instance
func NewCollectionHandlers(
	collectionServices crud_services.BaseCRUDService[entities.CollectionEntity],
	endpointServices crud_services.BaseCRUDService[entities.EndpointEntity],
	serviceContainer services.ServiceContainer,
) BaseCRUDHandlers[entities.CollectionEntity] {
	return &collectionHandlers{
		baseCRUDHandlers: NewBaseCRUDHandlers(collectionServices),
		endpointServices: endpointServices,
		serviceContainer: serviceContainer,
	}
}

func (h collectionHandlers) RegisterRoutes(e *echo.Group) {
	// Register standard CRUD routes
	h.baseCRUDHandlers.RegisterRoutes(e)

	// Register custom routes
	path := h.baseCRUDHandlers.GetPath()
	e.GET(fmt.Sprintf("/%s/:id/endpoints", path), h.getEndpoints)
	e.GET(fmt.Sprintf("/%s/:slug/endpoint_scenario", path), h.getEndpointScenario)
}

func (h collectionHandlers) GetPath() string {
	return h.baseCRUDHandlers.GetPath()
}

func (h collectionHandlers) GetServices() crud_services.BaseCRUDService[entities.CollectionEntity] {
	return h.baseCRUDHandlers.GetServices()
}

// getEndpoints returns all endpoints for a collection
func (h collectionHandlers) getEndpoints(ctx echo.Context) error {
	id, err := validateUUID(ctx)
	if err != nil {
		return respErr(ctx, http.StatusBadRequest, err)
	}

	query, err := extractPaginationQuery(ctx)
	if err != nil {
		return respErr(ctx, http.StatusBadRequest, err)
	}

	query.ShowAll = true
	query.Filters = make(map[string]interface{})
	query.Filters["collection_id = ?"] = id

	result, err := h.endpointServices.GetList(context.Background(), &query)
	if err != nil {
		return respErr(ctx, http.StatusInternalServerError, err)
	}

	return respJSON(ctx, http.StatusOK, result)
}

// getEndpointScenario returns endpoint scenarios for a collection by slug
func (h collectionHandlers) getEndpointScenario(ctx echo.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		return respErr(ctx, http.StatusBadRequest, fmt.Errorf("slug is required"))
	}

	endpointScenario, err := h.serviceContainer.GetEndpointScenarioService().GetEndpointScenario(context.Background(), slug)
	if err != nil {
		return respErr(ctx, http.StatusBadRequest, err)
	}

	return respJSON(ctx, http.StatusOK, endpointScenario)
}
