package handler

import (
	`context`
	`encoding/json`
	`fmt`
	`io`
	`net/http`
	`time`

	`berlin.allobank.com/tools/bonclay/internal/common/pagination`
	api_mock_services `berlin.allobank.com/tools/bonclay/internal/services/api-mock-services`
	`github.com/google/uuid`
	`github.com/labstack/echo/v4`
	`github.com/labstack/gommon/log`
)

func (h handlers) actMock(ctx echo.Context) error {
	method := ctx.Request().Method
	path := ctx.Param("path")
	collectionSlug := ctx.Param("collection")

	// Capture request body for history before it's consumed
	var reqBody string
	if ctx.Request().Body != nil {
		bodyBytes, _ := io.ReadAll(ctx.Request().Body)
		reqBody = string(bodyBytes)
		ctx.Request().Body = io.NopCloser(
			func() io.Reader {
				return &byteReader{data: bodyBytes}
			}(),
		)
	}

	// Capture request headers
	reqHeaders := make(map[string]string)
	for k, v := range ctx.Request().Header {
		if len(v) > 0 {
			reqHeaders[k] = v[0]
		}
	}

	startTime := time.Now()

	scenario, err := h.serviceContainer.GetMockService().MockAPI(ctx, collectionSlug, method, path)
	if err != nil {
		return ctx.String(404, err.Error())
	}

	// Extract response headers before diverging on proxy vs non-proxy
	respHeaders := make(map[string]string)
	if scenario.ProxyIsEnabled {
		for k, v := range scenario.ProxyResponseHeader {
			if len(v) > 0 {
				respHeaders[k] = v[0]
			}
		}
	} else {
		var parsedHeaders map[string]string
		if jsonErr := json.Unmarshal([]byte(scenario.Header), &parsedHeaders); jsonErr == nil {
			respHeaders = parsedHeaders
		}
	}

	// Estimate total duration: scenario lookup time + configured delay
	scenarioMs := time.Since(startTime).Milliseconds()
	var delayMs int64
	if scenario.Delay != nil && *scenario.Delay > 0 {
		delayMs = int64(*scenario.Delay) * 1000
	}
	durationMs := scenarioMs + delayMs

	// Record history asynchronously
	go func() {
		entry := api_mock_services.RequestHistory{
			ID:              uuid.New().String(),
			CollectionSlug:  collectionSlug,
			Method:          method,
			Path:            fmt.Sprintf("/%s", path),
			RequestHeaders:  reqHeaders,
			RequestBody:     reqBody,
			StatusCode:      scenario.StatusHeader,
			ResponseHeaders: respHeaders,
			DurationMs:      durationMs,
			ProxyUsed:       scenario.ProxyIsEnabled,
			ScenarioID:      scenario.EndpointID,
			Timestamp:       time.Now(),
		}
		if scenario.ProxyIsEnabled {
			entry.ResponseBody = string(scenario.ProxyResponseBody)
		} else {
			entry.ResponseBody = scenario.Body
		}
		_ = h.serviceContainer.GetHistoryService().Push(context.Background(), entry)
	}()

	if scenario.ProxyIsEnabled {
		log.Info("proxy run")
		for key, values := range scenario.ProxyResponseHeader {
			for _, value := range values {
				ctx.Response().Header().Add(key, value)
			}
		}

		if scenario.Delay != nil && *scenario.Delay > 0 {
			time.Sleep(time.Duration(*scenario.Delay) * time.Second)
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

// byteReader is a simple io.Reader over a byte slice (avoids bytes.NewReader import cycle)
type byteReader struct {
	data []byte
	pos  int
}

func (r *byteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// getHistory returns mock request history for a collection
func (h handlers) getHistory(ctx echo.Context) error {
	collectionSlug := ctx.Param("slug")
	if collectionSlug == "" {
		return respErr(ctx, http.StatusBadRequest, fmt.Errorf("slug is required"))
	}

	entries, err := h.serviceContainer.GetHistoryService().GetByCollection(context.Background(), collectionSlug, 100)
	if err != nil {
		return respErr(ctx, http.StatusInternalServerError, err)
	}

	return respJSON(ctx, http.StatusOK, entries)
}

// clearHistory deletes all history for a collection
func (h handlers) clearHistory(ctx echo.Context) error {
	collectionSlug := ctx.Param("slug")
	if collectionSlug == "" {
		return respErr(ctx, http.StatusBadRequest, fmt.Errorf("slug is required"))
	}

	if err := h.serviceContainer.GetHistoryService().Clear(context.Background(), collectionSlug); err != nil {
		return respErr(ctx, http.StatusInternalServerError, err)
	}

	return respJSON(ctx, http.StatusOK, map[string]string{"status": "cleared"})
}

// getSquadDetail returns squad with nested collections
func (h handlers) getSquadDetail(ctx echo.Context) error {
	c := context.Background()
	slug := ctx.Param("slug")

	if slug == "" {
		return respErr(ctx, http.StatusBadRequest, fmt.Errorf("slug is required"))
	}

	// Try by slug first, fall back to ID for squads that have no slug yet
	squadQuery := &pagination.ListQuery{
		ShowAll: true,
		Filters: map[string]interface{}{"slug = ?": slug},
	}
	squadResult, err := h.serviceContainer.GetCRUDServices().SquadServices.GetList(c, squadQuery)
	if err != nil || len(squadResult.List) == 0 {
		squadQuery = &pagination.ListQuery{
			ShowAll: true,
			Filters: map[string]interface{}{"id = ?": slug},
		}
		squadResult, err = h.serviceContainer.GetCRUDServices().SquadServices.GetList(c, squadQuery)
		if err != nil || len(squadResult.List) == 0 {
			return respErr(ctx, http.StatusNotFound, fmt.Errorf("squad not found"))
		}
	}
	squad := squadResult.List[0]

	collectionQuery := &pagination.ListQuery{
		ShowAll: true,
		Filters: map[string]interface{}{"squad_id = ?": squad.ID.String()},
		Search:  ctx.QueryParam("search"),
	}
	collectionsResult, err := h.serviceContainer.GetCRUDServices().CollectionServices.GetList(c, collectionQuery)
	if err != nil {
		return respErr(ctx, http.StatusInternalServerError, err)
	}

	response := map[string]interface{}{
		"id":          squad.ID,
		"name":        squad.Name,
		"slug":        squad.Slug,
		"desc":        squad.Desc,
		"collections": collectionsResult.List,
	}

	return respJSON(ctx, http.StatusOK, response)
}
