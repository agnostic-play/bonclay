package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	kong_services "berlin.allobank.com/tools/bonclay/internal/services/kong-services"
	"github.com/labstack/echo/v4"
)

type kongHandlers struct {
	kongService kong_services.KongService
}

func NewKongHandlers(kongService kong_services.KongService) *kongHandlers {
	return &kongHandlers{kongService: kongService}
}

func (h *kongHandlers) RegisterRoutes(e *echo.Group) {
	e.GET("/kong/services", h.listServices)
	e.PATCH("/kong/services/:id", h.updateService)
}

func (h *kongHandlers) updateService(c echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return respErr(c, http.StatusBadRequest, fmt.Errorf("service id is required"))
	}

	var payload kong_services.UpdateServicePayload
	if err := c.Bind(&payload); err != nil {
		return respErr(c, http.StatusBadRequest, err)
	}

	res, err := h.kongService.UpdateService(c.Request().Context(), id, payload)
	if err != nil {
		return respErr(c, http.StatusBadGateway, err)
	}

	return respJSON(c, http.StatusOK, res)
}

func (h *kongHandlers) listServices(c echo.Context) error {
	query := kong_services.ListServicesQuery{
		Offset: strings.TrimSpace(c.QueryParam("offset")),
		Tag:    strings.TrimSpace(c.QueryParam("tag")),
	}

	if sizeStr := strings.TrimSpace(c.QueryParam("size")); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 {
			query.Size = size
		}
	}

	res, err := h.kongService.ListServices(c.Request().Context(), query)
	if err != nil {
		return respErr(c, http.StatusBadGateway, err)
	}

	return respJSON(c, http.StatusOK, res)
}
