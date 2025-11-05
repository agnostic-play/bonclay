package handler

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h handlers) routesViews() {
	h.server.GET("/share/diagram/:url", h.viewShareDiagram)
	h.server.GET("/squad/:slug", h.viewShowSquad)
	h.server.GET("/squad/:slug/collection/:collectionSlug", h.viewShowCollection)

	//h.server.Static("/v2/assets", "resources/views/spa/dist/assets")
	//h.server.File("/v2", "resources/views/spa/dist/index.html")
	//h.server.File("/v2/*", "resources/views/spa/dist/index.html")

	h.server.Static("/v2/assets", "resources/views/bonclay/public/bonclay/assets")
	h.server.File("/v2", "resources/views/bonclay/public/bonclay/index.html")
	h.server.File("/v2/*", "resources/views/bonclay/public/bonclay/index.html")
}

func (h handlers) viewShowCollection(ctx echo.Context) error {
	var collection interface{}
	var endpoints interface{}

	collection = nil
	endpoints = nil

	c := context.Background()

	squad, err := h.serviceContainer.GetSquad(c, ctx.Param("slug"))
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, h.baseURL)
		return nil
	}

	collection, err = h.serviceContainer.GetCollection(c, ctx.Param("collectionSlug"))
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, h.baseURL+"/squad/"+squad.Slug)
		return nil
	}

	endpointScenario, _ := h.serviceContainer.GetEndpointScenario(c, ctx.Param("collectionSlug"))
	if len(endpointScenario) > 0 {
		endpoints = endpointScenario
	}

	return h.render(ctx, "squad_show", map[string]interface{}{
		"squad":      squad,
		"collection": collection,
		"endpoints":  endpoints,
	})
}

func (h handlers) viewShowSquad(ctx echo.Context) error {
	var collection interface{}
	collection = nil
	c := context.Background()

	resp, err := h.serviceContainer.GetSquad(c, ctx.Param("slug"))
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, h.baseURL)
		return nil
	}

	if collectionSlug := ctx.Param("collectionSlug"); collectionSlug != "" {
		collection, _ = h.serviceContainer.GetCollection(c, collectionSlug)
	}

	return h.render(ctx, "squad_show", map[string]interface{}{
		"squad":      resp,
		"collection": collection,
	})
}

func (h handlers) viewShareDiagram(ctx echo.Context) error {
	c := context.Background()

	url := ctx.Param("url")
	if url == "" {
		return ctx.String(http.StatusBadRequest, "url is required")
	}

	decodeString, err := base64.StdEncoding.DecodeString(url)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "url is required")
	}

	data, err := h.serviceContainer.GetCRUDServices().DiagramServices.Get(c, string(decodeString))
	if err != nil {
		log.Println("error getting diagram service:", err)
		// Consider returning 404 here instead of 400
		return ctx.String(http.StatusNotFound, "diagram not found")
	}

	// ✅ Pass only what the template needs; avoids field-path mistakes.

	if err := h.render(ctx, "shared_mermaid_diagram", data); err != nil {
		log.Println("error rendering diagram:", err)
		return ctx.String(http.StatusNotFound, "diagram not found")
	}

	return nil
}
