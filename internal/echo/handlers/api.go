package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/agnostic-play/ditoo/internal/services"
	"github.com/labstack/echo/v4"
)

func (h handlers) routesApi() {
	h.server.POST("/api/squad/create", h.actCreateSquad)
	h.server.PUT("/api/squad/update/:id", h.actUpdateSquad)
	h.server.DELETE("/api/squad/delete/:id", h.actDeleteSquad)
	h.server.GET("/api/squad", h.getListSquad)

	h.server.POST("/api/collection/create", h.actCreateCollection)
	h.server.PUT("/api/collection/update/:id", h.actUpdateCollection)
	h.server.DELETE("/api/collection/delete/:id", h.actDeleteCollection)
	h.server.GET("/api/collection/:slug/endpoint", h.actGetEndpoint)
	h.server.GET("/api/collection/endpoint_scenario/:collectionSlug", h.getCollectionEndpoint)

	h.server.GET("/api/collection/custom_variable/list/:collectionSlug", h.getCustomVariable)
	h.server.POST("/api/collection/custom_variable/create", h.actCreateCustomVariable)

	h.server.POST("/api/endpoint/create", h.actCreateEndpoint)
	h.server.PUT("/api/endpoint/update/:id", h.actUpdateEndpoint)
	h.server.DELETE("/api/endpoint/delete/:id", h.actDeleteEndpoint)

	h.server.POST("/api/scenario/create", h.actCreateScenario)
	h.server.PUT("/api/scenario/update/:id", h.actUpdateScenario)
	h.server.DELETE("/api/scenario/delete/:id", h.actDeleteScenario)

	h.server.POST("/api/set/active_scenario", h.actSetActiveScenario)
	h.server.POST("/api/v2/tools/config-encryption/encrypt", h.actEncryptConfig)
	h.server.POST("/api/v2/tools/config-encryption/decrypt", h.actDecryptConfig)

	h.server.Any("/mock/:collection/:path", h.actMock)

}

func (h handlers) actMock(ctx echo.Context) error {
	method := ctx.Request().Method
	path := ctx.Param("path")
	collectionSlug := ctx.Param("collection")

	scenario, err := h.serviceContainer.MockApi(context.Background(), collectionSlug, method, path)
	if err != nil {
		return ctx.String(404, err.Error())
	}

	var header map[string]string

	err = json.Unmarshal([]byte(scenario.Header), &header)
	if err != nil {
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

func (h handlers) getListSquad(ctx echo.Context) error {
	name := ctx.QueryParam("name")

	resp, err := h.serviceContainer.GetListSquad(context.Background(), name)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
}

func (h handlers) actCreateSquad(ctx echo.Context) error {
	var req services.SquadEntityReq

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.CreateOrUpdateSquad(context.Background(), "", req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
}

func (h handlers) actUpdateSquad(ctx echo.Context) error {
	var req services.SquadEntityReq

	id, err := h.validateUUID(ctx)
	if err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.CreateOrUpdateSquad(context.Background(), id, req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
}

func (h handlers) actDeleteSquad(ctx echo.Context) error {

	id, err := h.validateUUID(ctx)
	if err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.DeleteSquad(context.Background(), id)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}
	return h.json(ctx, 200, resp)

}

func (h handlers) actGetEndpoint(ctx echo.Context) error {

	resp, err := h.serviceContainer.GetEndpointScenario(context.Background(), ctx.Param("slug"))
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
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

func (h handlers) actUpdateCollection(ctx echo.Context) error {
	var req services.CollectionEntityReq

	id, err := h.validateUUID(ctx)
	if err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.CreateOrUpdateCollection(context.Background(), id, req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
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

func (h handlers) getCollectionEndpoint(ctx echo.Context) error {
	c := context.Background()

	endpointScenario, err := h.serviceContainer.GetEndpointScenario(c, ctx.Param("collectionSlug"))
	if err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	return h.json(ctx, 200, endpointScenario)
}

func (h handlers) getCustomVariable(ctx echo.Context) error {
	c := context.Background()

	customVariable, err := h.serviceContainer.GetCustomVariable(c, ctx.Param("collectionSlug"))
	if err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	return h.json(ctx, 200, customVariable)
}

func (h handlers) actCreateCustomVariable(ctx echo.Context) error {
	var req services.CustomVariableReq

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.CreateOrUpdateCustomVariable(context.Background(), "", req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
}

func (h handlers) actCreateEndpoint(ctx echo.Context) error {
	var req services.EndpointEntityReq

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.CreateOrUpdateEndpoint(context.Background(), "", req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
}

func (h handlers) actUpdateEndpoint(ctx echo.Context) error {
	var req services.EndpointEntityReq

	id, err := h.validateUUID(ctx)
	if err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.CreateOrUpdateEndpoint(context.Background(), id, req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
}

func (h handlers) actDeleteEndpoint(ctx echo.Context) error {

	id, err := h.validateUUID(ctx)
	if err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	err = h.serviceContainer.DeleteEndpoint(context.Background(), id)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}
	return h.json(ctx, 200, nil)
}

func (h handlers) actCreateScenario(ctx echo.Context) error {
	var req services.ScenarioEntityReq

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.CreateOrUpdateScenario(context.Background(), "", req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
}

func (h handlers) actUpdateScenario(ctx echo.Context) error {
	var req services.ScenarioEntityReq

	id, err := h.validateUUID(ctx)
	if err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.CreateOrUpdateScenario(context.Background(), id, req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
}

func (h handlers) actDeleteScenario(ctx echo.Context) error {

	id, err := h.validateUUID(ctx)
	if err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	err = h.serviceContainer.DeleteScenario(context.Background(), id)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}
	return h.json(ctx, 200, nil)
}

func (h handlers) actSetActiveScenario(ctx echo.Context) error {
	var req services.SetActiveScenarioEntityReq

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	err := h.serviceContainer.SetActiveResponse(context.Background(), req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, req)
}

func (h handlers) actEncryptConfig(ctx echo.Context) error {
	var req services.EncryptionToolsEntity

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.EncryptConfig(context.Background(), req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
}

func (h handlers) actDecryptConfig(ctx echo.Context) error {
	var req services.EncryptionToolsEntity

	if err := h.validateRequest(ctx, &req); err != nil {
		return h.errorJson(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.DecryptConfig(context.Background(), req)
	if err != nil {
		return h.errorJson(ctx, http.StatusInternalServerError, err)
	}

	return h.json(ctx, 200, resp)
}
