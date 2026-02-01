package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/agnostic-play/ditoo/internal/entities"
	"github.com/agnostic-play/ditoo/internal/services"
	crud_services "github.com/agnostic-play/ditoo/internal/services/crud-services"
	"github.com/labstack/echo/v4"
)

// CollectionRoutes extends BaseCRUDHandlers with collection-specific endpoints
type CollectionRoutes interface {
	BaseCRUDHandlers[entities.CollectionEntity]
}

type collectionRoutes struct {
	baseCRUDRoutes   BaseCRUDHandlers[entities.CollectionEntity]
	endpointServices crud_services.BaseCRUDService[entities.EndpointEntity]
	serviceContainer services.ServiceContainer
}

// NewCollectionRoutes creates a new CollectionRoutes instance
func NewCollectionRoutes(
	collectionServices crud_services.BaseCRUDService[entities.CollectionEntity],
	endpointServices crud_services.BaseCRUDService[entities.EndpointEntity],
	serviceContainer services.ServiceContainer,
) BaseCRUDHandlers[entities.CollectionEntity] {
	return &collectionRoutes{
		baseCRUDRoutes:   NewBaseCRUDHandlers(collectionServices),
		endpointServices: endpointServices,
		serviceContainer: serviceContainer,
	}
}

func (h collectionRoutes) RegisterRoutes(e *echo.Group) {
	// Register standard CRUD routes
	h.baseCRUDRoutes.RegisterRoutes(e)

	// Register custom routes
	path := h.baseCRUDRoutes.GetPath()
	e.GET(fmt.Sprintf("/%s/:id/endpoints", path), h.getEndpoints)
	e.GET(fmt.Sprintf("/%s/:slug/endpoint_scenario", path), h.getEndpointScenario)
}

func (h collectionRoutes) GetPath() string {
	return h.baseCRUDRoutes.GetPath()
}

func (h collectionRoutes) GetServices() crud_services.BaseCRUDService[entities.CollectionEntity] {
	return h.baseCRUDRoutes.GetServices()
}

// getEndpoints returns all endpoints for a collection
func (h collectionRoutes) getEndpoints(ctx echo.Context) error {
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
func (h collectionRoutes) getEndpointScenario(ctx echo.Context) error {
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
