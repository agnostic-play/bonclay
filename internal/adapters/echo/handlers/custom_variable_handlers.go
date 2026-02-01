package handler

import (
	`context`
	`net/http`

	api_mock_services `github.com/agnostic-play/ditoo/internal/services/api-mock-services`
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
