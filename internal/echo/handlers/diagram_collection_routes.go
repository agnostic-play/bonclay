package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/agnostic-play/ditoo/internal/entities"
	crud_services "github.com/agnostic-play/ditoo/internal/services/crud-services"
	"github.com/labstack/echo/v4"
)

type DiagramCollectionRoutes interface {
	BaseCRUDRoutes[entities.DiagramCollectionEntity]
}

type diagramCollectionRoutes struct {
	diagramCollectionRoutes BaseCRUDRoutes[entities.DiagramCollectionEntity]
	diagramServices         crud_services.BaseCRUDService[entities.DiagramEntity]
}

func NewDiagramCollectionRoutes(diagramCollectionServices crud_services.BaseCRUDService[entities.DiagramCollectionEntity], diagramService crud_services.BaseCRUDService[entities.DiagramEntity]) BaseCRUDRoutes[entities.DiagramCollectionEntity] {
	return &diagramCollectionRoutes{
		diagramCollectionRoutes: NewBaseCRUDRoutes(diagramCollectionServices),
		diagramServices:         diagramService,
	}
}

func (h diagramCollectionRoutes) RegisterRoutes(e *echo.Group) {
	h.diagramCollectionRoutes.RegisterRoutes(e)
	e.GET(fmt.Sprintf("/%s/:id/diagrams", h.diagramCollectionRoutes.GetPath()), h.getDiagrams)
}

func (h diagramCollectionRoutes) GetPath() string {
	return h.diagramCollectionRoutes.GetPath()
}

func (h diagramCollectionRoutes) GetServices() crud_services.BaseCRUDService[entities.DiagramCollectionEntity] {
	return h.GetServices()
}

func (h diagramCollectionRoutes) getDiagrams(ctx echo.Context) error {
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
