package handler

import (
	`context`
	`encoding/json`
	`fmt`
	`net/http`
	`time`

	`github.com/agnostic-play/ditoo/internal/common/pagination`
	`github.com/labstack/echo/v4`
	`github.com/labstack/gommon/log`
)

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

// getSquadDetail returns squad with nested collections
func (h handlers) getSquadDetail(ctx echo.Context) error {
	c := context.Background()
	slug := ctx.Param("slug")

	if slug == "" {
		return respErr(ctx, http.StatusBadRequest, fmt.Errorf("slug is required"))
	}

	// Get squad by slug
	squadQuery := &pagination.ListQuery{
		ShowAll: true,
		Filters: map[string]interface{}{"slug = ?": slug},
	}
	squadResult, err := h.serviceContainer.GetCRUDServices().SquadServices.GetList(c, squadQuery)
	if err != nil || len(squadResult.List) == 0 {
		return respErr(ctx, http.StatusNotFound, fmt.Errorf("squad not found"))
	}
	squad := squadResult.List[0]

	// Get collections for this squad
	collectionQuery := &pagination.ListQuery{
		ShowAll: true,
		Filters: map[string]interface{}{"squad_id = ?": squad.ID.String()},
	}
	collectionsResult, err := h.serviceContainer.GetCRUDServices().CollectionServices.GetList(c, collectionQuery)
	if err != nil {
		return respErr(ctx, http.StatusInternalServerError, err)
	}

	// Build response with squad and collections
	response := map[string]interface{}{
		"id":          squad.ID,
		"name":        squad.Name,
		"slug":        squad.Slug,
		"desc":        squad.Desc,
		"collections": collectionsResult.List,
	}

	return respJSON(ctx, http.StatusOK, response)
}
