package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/agnostic-play/ditoo/internal/common/pagination"
	crud_services "github.com/agnostic-play/ditoo/internal/services/crud-services"
	"github.com/labstack/echo/v4"
)

type BaseCRUDRoutes[T any] interface {
	RegisterRoutes(e *echo.Group)
	GetPath() string
	GetServices() crud_services.BaseCRUDService[T]
}

type baseCRUDRoutes[T any] struct {
	services crud_services.BaseCRUDService[T]
	entity   T
	path     string
	title    string
}

func NewBaseCRUDRoutes[T any](services crud_services.BaseCRUDService[T]) BaseCRUDRoutes[T] {
	ent := services.GetEntity()
	return &baseCRUDRoutes[T]{
		services: services,
		title:    ent.GetTitle(),
		path:     ent.GetTableName(),
		entity:   *ent.GetEntity(),
	}
}

func (h baseCRUDRoutes[T]) RegisterRoutes(e *echo.Group) {
	e.GET(fmt.Sprintf("/%s/all", h.path), h.getAll)
	e.GET(fmt.Sprintf("/%s/list", h.path), h.getPagination)
	e.POST(fmt.Sprintf("/%s/create", h.path), h.actCreate)
	e.GET(fmt.Sprintf("/%s/:id/show", h.path), h.getShow)
	e.PATCH(fmt.Sprintf("/%s/:id/update", h.path), h.actUpdate)
	e.DELETE(fmt.Sprintf("/%s/:id/remove", h.path), h.actRemove)
}

func (h baseCRUDRoutes[T]) GetServices() crud_services.BaseCRUDService[T] {
	return h.services
}

func (h baseCRUDRoutes[T]) GetPath() string {
	return h.path
}

// GET {path}/all?search=...&sort_by=...  (uses ListQuery but returns all that match)
func (h baseCRUDRoutes[T]) getAll(c echo.Context) error {
	q, _ := extractPaginationQuery(c)
	q.ShowAll = true

	res, err := h.services.GetList(context.Background(), &q)
	if err != nil {
		return respErrJSON(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, res)
}

// GET {path}/list?page=..&page_size=..
func (h baseCRUDRoutes[T]) getPagination(c echo.Context) error {
	q, err := extractPaginationQuery(c)
	if err != nil {
		return respErrJSON(c, http.StatusBadRequest, err)
	}

	// Normalize limit/offset if only page/page_size given
	if q.Limit == 0 {
		q.Limit = q.PageSize
	}
	if q.Offset == 0 {
		q.Offset = (q.Page - 1) * q.PageSize
	}
	res, err := h.services.GetList(context.Background(), &q)
	if err != nil {
		return respErrJSON(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, res)
}

// POST {path}/create
func (h baseCRUDRoutes[T]) actCreate(c echo.Context) error {
	var req T
	if err := validateRequest(c, &req); err != nil {
		return respErrJSON(c, http.StatusBadRequest, err)
	}
	created, err := h.services.Create(context.Background(), req)
	if err != nil {
		return respErrJSON(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, created)
}

// GET {path}/:id/show
func (h baseCRUDRoutes[T]) getShow(c echo.Context) error {
	id, err := validateUUID(c)
	if err != nil {
		return respErrJSON(c, http.StatusBadRequest, err)
	}
	item, err := h.services.Get(context.Background(), id)
	if err != nil {
		return respErrJSON(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, item)
}

// PATCH {path}/:id/update
func (h baseCRUDRoutes[T]) actUpdate(c echo.Context) error {
	id, err := validateUUID(c)
	if err != nil {
		return respErrJSON(c, http.StatusBadRequest, err)
	}
	var req T
	if err := validateRequest(c, &req); err != nil {
		return respErrJSON(c, http.StatusBadRequest, err)
	}
	updated, err := h.services.Patch(context.Background(), id, req)
	if err != nil {
		return respErrJSON(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, updated)
}

// DELETE {path}/:id/remove
func (h baseCRUDRoutes[T]) actRemove(c echo.Context) error {
	id, err := validateUUID(c)
	if err != nil {
		return respErrJSON(c, http.StatusBadRequest, err)
	}
	if err := h.services.Delete(context.Background(), id); err != nil {
		return respErrJSON(c, http.StatusInternalServerError, err)
	}
	return respJSON(c, http.StatusOK, echo.Map{"status": "ok"})
}

func extractPaginationQuery(c echo.Context) (pagination.ListQuery, error) {
	var query pagination.ListQuery

	// Page (default 1)
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			query.Page = page
		}
	}
	if query.Page == 0 {
		query.Page = 1
	}

	// PageSize (default 10)
	if sizeStr := c.QueryParam("page_size"); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 {
			query.PageSize = size
		}
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	// Limit & Offset (optional overrides)
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			query.Limit = limit
		}
	}
	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			query.Offset = offset
		}
	}

	// Search
	query.Search = c.QueryParam("search")

	// Sorting
	query.SortBy = c.QueryParam("sort_by")
	query.SortOrder = c.QueryParam("sort_order")

	return query, nil
}
