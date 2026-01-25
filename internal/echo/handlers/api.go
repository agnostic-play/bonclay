package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	api_mock_services "github.com/agnostic-play/ditoo/internal/services/api-mock-services"
	common_services "github.com/agnostic-play/ditoo/internal/services/common-services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (h handlers) routesApi() {
	// Legacy routes (not yet migrated to v2)
	h.server.GET("/api/collection/custom_variable/list/:collectionSlug", h.getCustomVariable)
	h.server.POST("/api/collection/custom_variable/create", h.actCreateCustomVariable)

	// V2 API - All routes use base CRUD pattern
	v2 := h.server.Group("/api/v2")
	h.registerV2Routes(v2)

	// Mock API
	h.server.Any("/mock/:collection/:path", h.actMock)
}

func (h handlers) registerV2Routes(v2 *echo.Group) {
	crudServices := h.serviceContainer.GetCRUDServices()

	// Tools
	v2.POST("/tools/config-encryption/encrypt", h.actEncryptConfig)
	v2.POST("/tools/config-encryption/decrypt", h.actDecryptConfig)

	// Diagrams
	NewDiagramCollectionRoutes(crudServices.DiagramCollectionServices, crudServices.DiagramServices).RegisterRoutes(v2)
	NewBaseCRUDRoutes(crudServices.DiagramServices).RegisterRoutes(v2)

	// Squad, Collection, Endpoint, Scenario
	NewBaseCRUDRoutes(crudServices.SquadServices).RegisterRoutes(v2)
	NewCollectionRoutes(crudServices.CollectionServices, crudServices.EndpointServices, h.serviceContainer).RegisterRoutes(v2)
	NewBaseCRUDRoutes(crudServices.EndpointServices).RegisterRoutes(v2)
	NewScenarioRoutes(crudServices.ScenarioServices, h.serviceContainer).RegisterRoutes(v2)
}

func (h handlers) actMock(ctx echo.Context) error {
	method := ctx.Request().Method
	path := ctx.Param("path")
	collectionSlug := ctx.Param("collection")

	scenario, err := h.serviceContainer.GetMockService().MockAPI(ctx, collectionSlug, method, path)
	if err != nil {
		return ctx.String(404, err.Error())
	}

	if scenario.ProxyIsEnabled {
		log.Info("proxy run")
		for key, values := range scenario.ProxyResponseHeader {
			for _, value := range values {
				ctx.Response().Header().Add(key, value)
			}
		}

		ctx.Response().WriteHeader(scenario.StatusHeader)
		_, err = ctx.Response().Write(scenario.ProxyResponseBody)
		if err != nil {
			return ctx.String(500, err.Error())
		}

		return nil
	}

	var header map[string]string

	if err = json.Unmarshal([]byte(scenario.Header), &header); err != nil {
		return ctx.String(500, err.Error())
	}

	for index, val := range header {
		ctx.Response().Header().Set(index, val)
	}

	ctx.Response().WriteHeader(scenario.StatusHeader)

	if scenario.Delay != nil && *scenario.Delay > 0 {
		time.Sleep(time.Duration(*scenario.Delay) * time.Second)
	}

	_, err = ctx.Response().Write([]byte(scenario.Body))
	if err != nil {
		return ctx.String(500, err.Error())
	}
	return nil

}

// Legacy custom variable handlers
func (h handlers) getCustomVariable(ctx echo.Context) error {
	c := context.Background()

	customVariable, err := h.serviceContainer.GetCustomVariableService().GetCustomVariable(c, ctx.Param("collectionSlug"))
	if err != nil {
		return respErrJSON(ctx, http.StatusBadRequest, err)
	}

	return respJSON(ctx, 200, customVariable)
}

func (h handlers) actCreateCustomVariable(ctx echo.Context) error {
	var req api_mock_services.CustomVariableReq

	if err := validateRequest(ctx, &req); err != nil {
		return respErrJSON(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.GetCustomVariableService().CreateOrUpdateCustomVariable(context.Background(), "", req)
	if err != nil {
		return respErrJSON(ctx, http.StatusInternalServerError, err)
	}

	return respJSON(ctx, 200, resp)
}

func (h handlers) actEncryptConfig(ctx echo.Context) error {
	var req common_services.EncryptionReq

	if err := validateRequest(ctx, &req); err != nil {
		return respErrJSON(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.GetEncryptionService().EncryptConfig(context.Background(), req)
	if err != nil {
		return respErrJSON(ctx, http.StatusInternalServerError, err)
	}

	return respJSON(ctx, 200, resp)
}

func (h handlers) actDecryptConfig(ctx echo.Context) error {
	var req common_services.EncryptionReq

	if err := validateRequest(ctx, &req); err != nil {
		return respErrJSON(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.GetEncryptionService().DecryptConfig(context.Background(), req)
	if err != nil {
		return respErrJSON(ctx, http.StatusInternalServerError, err)
	}

	return respJSON(ctx, 200, resp)
}
