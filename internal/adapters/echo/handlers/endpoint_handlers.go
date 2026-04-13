package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/agnostic-play/ditoo/internal/common/pagination"
	"github.com/agnostic-play/ditoo/internal/entities"
	crud_services "github.com/agnostic-play/ditoo/internal/services/crud-services"
	"github.com/labstack/echo/v4"
)

type endpointRoutes struct {
	baseCRUDRoutes BaseCRUDHandlers[entities.EndpointEntity]
}

// NewEndpointHandlers wraps the base CRUD handler and adds duplicate validation on create.
func NewEndpointHandlers(services crud_services.BaseCRUDService[entities.EndpointEntity]) BaseCRUDHandlers[entities.EndpointEntity] {
	return &endpointRoutes{
		baseCRUDRoutes: NewBaseCRUDHandlers(services),
	}
}

func (h *endpointRoutes) RegisterRoutes(e *echo.Group) {
	path := h.baseCRUDRoutes.GetPath()

	// Override create with duplicate validation; delegate everything else to base.
	e.GET(fmt.Sprintf("/%s/all", path), h.proxyGetAll)
	e.GET(fmt.Sprintf("/%s/list", path), h.proxyGetList)
	e.POST(fmt.Sprintf("/%s/create", path), h.actCreate)
	e.GET(fmt.Sprintf("/%s/:id/show", path), h.proxyShow)
	e.PATCH(fmt.Sprintf("/%s/:id/update", path), h.proxyUpdate)
	e.DELETE(fmt.Sprintf("/%s/:id/remove", path), h.proxyDelete)
}

func (h *endpointRoutes) GetPath() string { return h.baseCRUDRoutes.GetPath() }
func (h *endpointRoutes) GetServices() crud_services.BaseCRUDService[entities.EndpointEntity] {
	return h.baseCRUDRoutes.GetServices()
}

// actCreate validates no duplicate (collection_id + method + path) before creating.
func (h *endpointRoutes) actCreate(c echo.Context) error {
	var req entities.EndpointEntity
	if err := validateRequest(c, &req); err != nil {
		return respErr(c, http.StatusBadRequest, err)
	}

	// Normalise for comparison
	req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	req.Path = strings.TrimSpace(req.Path)

	query := &pagination.ListQuery{
		ShowAll: true,
		Filters: map[string]interface{}{
			"collection_id = ?": req.CollectionID,
			"method = ?":        req.Method,
			"path = ?":          req.Path,
		},
	}
	existing, err := h.baseCRUDRoutes.GetServices().GetList(context.Background(), query)
	if err != nil {
		return respErr(c, http.StatusInternalServerError, err)
	}
	if existing != nil && len(existing.List) > 0 {
		return respErr(c, http.StatusConflict,
			fmt.Errorf("endpoint %s %s already exists in this collection", req.Method, req.Path))
	}

	created, err := h.baseCRUDRoutes.GetServices().Create(context.Background(), req)
	if err != nil {
		return respErr(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, created)
}

// proxy helpers — delegate to the base handler's registered logic via service calls directly.

func (h *endpointRoutes) proxyGetAll(c echo.Context) error {
	q, _ := extractPaginationQuery(c)
	q.ShowAll = true
	res, err := h.baseCRUDRoutes.GetServices().GetList(context.Background(), &q)
	if err != nil {
		return respErr(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, res)
}

func (h *endpointRoutes) proxyGetList(c echo.Context) error {
	q, err := extractPaginationQuery(c)
	if err != nil {
		return respErr(c, http.StatusBadRequest, err)
	}
	if q.Limit == 0 {
		q.Limit = q.PageSize
	}
	if q.Offset == 0 {
		q.Offset = (q.Page - 1) * q.PageSize
	}
	res, err := h.baseCRUDRoutes.GetServices().GetList(context.Background(), &q)
	if err != nil {
		return respErr(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, res)
}

func (h *endpointRoutes) proxyShow(c echo.Context) error {
	id, err := validateUUID(c)
	if err != nil {
		return respErr(c, http.StatusBadRequest, err)
	}
	item, err := h.baseCRUDRoutes.GetServices().Get(context.Background(), id)
	if err != nil {
		return respErr(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, item)
}

func (h *endpointRoutes) proxyUpdate(c echo.Context) error {
	id, err := validateUUID(c)
	if err != nil {
		return respErr(c, http.StatusBadRequest, err)
	}
	var req entities.EndpointEntity
	if err := validateRequest(c, &req); err != nil {
		return respErr(c, http.StatusBadRequest, err)
	}
	updated, err := h.baseCRUDRoutes.GetServices().Patch(context.Background(), id, req)
	if err != nil {
		return respErr(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, updated)
}

func (h *endpointRoutes) proxyDelete(c echo.Context) error {
	id, err := validateUUID(c)
	if err != nil {
		return respErr(c, http.StatusBadRequest, err)
	}
	if err := h.baseCRUDRoutes.GetServices().Delete(context.Background(), id); err != nil {
		return respErr(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, echo.Map{"status": "ok"})
}
