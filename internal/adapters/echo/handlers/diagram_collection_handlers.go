package handler

import (
	"context"
	"fmt"
	"net/http"

	"berlin.allobank.com/tools/bonclay/internal/entities"
	crud_services "berlin.allobank.com/tools/bonclay/internal/services/crud-services"
	"github.com/labstack/echo/v4"
)

type DiagramCollectionHandlers interface {
	BaseCRUDHandlers[entities.DiagramCollectionEntity]
}

type diagramCollectionHandlers struct {
	diagramCollectionRoutes BaseCRUDHandlers[entities.DiagramCollectionEntity]
	diagramServices         crud_services.BaseCRUDService[entities.DiagramEntity]
}

func NewDiagramCollectionHandlers(diagramCollectionServices crud_services.BaseCRUDService[entities.DiagramCollectionEntity], diagramService crud_services.BaseCRUDService[entities.DiagramEntity]) BaseCRUDHandlers[entities.DiagramCollectionEntity] {
	return &diagramCollectionHandlers{
		diagramCollectionRoutes: NewBaseCRUDHandlers(diagramCollectionServices),
		diagramServices:         diagramService,
	}
}

func (h diagramCollectionHandlers) RegisterRoutes(e *echo.Group) {
	h.diagramCollectionRoutes.RegisterRoutes(e)
	e.GET(fmt.Sprintf("/%s/:id/diagrams", h.diagramCollectionRoutes.GetPath()), h.getDiagrams)
}

func (h diagramCollectionHandlers) GetPath() string {
	return h.diagramCollectionRoutes.GetPath()
}

func (h diagramCollectionHandlers) GetServices() crud_services.BaseCRUDService[entities.DiagramCollectionEntity] {
	return h.GetServices()
}

func (h diagramCollectionHandlers) getDiagrams(ctx echo.Context) error {
	id, err := validateUUID(ctx)
	if err != nil {
		return respErrJSON(ctx, http.StatusBadRequest, err)
	}

	query, err := extractPaginationQuery(ctx)
	if err != nil {
		return respErrJSON(ctx, http.StatusBadRequest, err)
	}

	query.ShowAll = true
	query.Filters = make(map[string]interface{})
	query.Filters["collection_id = ?"] = id

	result, err := h.diagramServices.GetList(context.Background(), &query)
	if err != nil {
		return respErrJSON(ctx, http.StatusInternalServerError, err)
	}

	return respJSON(ctx, http.StatusOK, result)
}
