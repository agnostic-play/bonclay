package handler

import (
	`context`
	`fmt`
	`net/http`

	api_mock_services `berlin.allobank.com/tools/bonclay/internal/services/api-mock-services`
	`github.com/labstack/echo/v4`
)

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

// v2 handlers

// GET /api/v2/collections/:slug/variables
func (h handlers) listCollectionVariables(ctx echo.Context) error {
	vars, err := h.serviceContainer.GetCustomVariableService().ListVariables(context.Background(), ctx.Param("slug"))
	if err != nil {
		return respErrJSON(ctx, http.StatusBadRequest, err)
	}
	return respJSON(ctx, http.StatusOK, vars)
}

// POST /api/v2/collections/:slug/variables  body: {key, value, collection_id}
func (h handlers) createCollectionVariable(ctx echo.Context) error {
	collectionID := ctx.QueryParam("collection_id")
	if collectionID == "" {
		return respErrJSON(ctx, http.StatusBadRequest, fmt.Errorf("collection_id query param is required"))
	}

	var req api_mock_services.CollectionVariableReq
	if err := validateRequest(ctx, &req); err != nil {
		return respErrJSON(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.GetCustomVariableService().CreateVariable(context.Background(), collectionID, req.Key, req.Value)
	if err != nil {
		return respErrJSON(ctx, http.StatusInternalServerError, err)
	}
	return respJSON(ctx, http.StatusOK, resp)
}

// PATCH /api/v2/collections/variables/:id  body: {key, value}
func (h handlers) updateCollectionVariable(ctx echo.Context) error {
	id, err := validateUUID(ctx)
	if err != nil {
		return respErrJSON(ctx, http.StatusBadRequest, err)
	}

	var req api_mock_services.CollectionVariableReq
	if err := validateRequest(ctx, &req); err != nil {
		return respErrJSON(ctx, http.StatusBadRequest, err)
	}

	resp, err := h.serviceContainer.GetCustomVariableService().UpdateVariable(context.Background(), id, req.Key, req.Value)
	if err != nil {
		return respErrJSON(ctx, http.StatusInternalServerError, err)
	}
	return respJSON(ctx, http.StatusOK, resp)
}

// DELETE /api/v2/collections/variables/:id
func (h handlers) deleteCollectionVariable(ctx echo.Context) error {
	id, err := validateUUID(ctx)
	if err != nil {
		return respErrJSON(ctx, http.StatusBadRequest, err)
	}

	if err := h.serviceContainer.GetCustomVariableService().DeleteVariable(context.Background(), id); err != nil {
		return respErrJSON(ctx, http.StatusInternalServerError, err)
	}
	return respJSON(ctx, http.StatusOK, echo.Map{"status": "ok"})
}
