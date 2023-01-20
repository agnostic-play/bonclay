package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"gitlab.linkaja.com/be/ditto/internal/services"
	"net/http"
)

func (h handlers) routesCollection() {
	h.server.POST("/api/collection/create", h.actCreateCollection)
	h.server.DELETE("/api/collection/delete/:id", h.actDeleteCollection)
}

func (h handlers) actCreateCollection(ctx echo.Context) error {
	var req services.CollectionEntityReq

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.CreateOrUpdateCollection(context.Background(), "", req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
}

func (h handlers) actUpdateCollection(c echo.Context) error {
	return nil
}

func (h handlers) actDeleteCollection(ctx echo.Context) error {

	id, err := h.validateUUID(ctx)
	if err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.DeleteCollection(context.Background(), id)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}
	return h.json(ctx, 200, resp)
}
