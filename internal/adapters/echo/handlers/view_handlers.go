package handler

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"

	"berlin.allobank.com/tools/bonclay/internal/common/pagination"
	"github.com/labstack/echo/v4"
)

func (h handlers) viewShowCollection(ctx echo.Context) error {
	var collection interface{}
	var endpoints interface{}

	collection = nil
	endpoints = nil

	c := context.Background()

	// Get squad by slug
	squadQuery := &pagination.ListQuery{ShowAll: true, Filters: map[string]interface{}{"slug = ?": ctx.Param("slug")}}
	squadResult, err := h.serviceContainer.GetCRUDServices().SquadServices.GetList(c, squadQuery)
	if err != nil || len(squadResult.List) == 0 {
		ctx.Redirect(http.StatusTemporaryRedirect, h.baseURL)
		return nil
	}
	squad := squadResult.List[0]

	// Get collection by slug
	collQuery := &pagination.ListQuery{ShowAll: true, Filters: map[string]interface{}{"slug = ?": ctx.Param("collectionSlug")}}
	collResult, err := h.serviceContainer.GetCRUDServices().CollectionServices.GetList(c, collQuery)
	if err != nil || len(collResult.List) == 0 {
		ctx.Redirect(http.StatusTemporaryRedirect, h.baseURL+"/squad/"+squad.Slug)
		return nil
	}
	collection = collResult.List[0]

	endpointScenario, _ := h.serviceContainer.GetEndpointScenarioService().GetEndpointScenario(c, ctx.Param("collectionSlug"))
	if len(endpointScenario.Endpoints) > 0 {
		endpoints = endpointScenario.Endpoints
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

	// Get squad by slug
	squadQuery := &pagination.ListQuery{ShowAll: true, Filters: map[string]interface{}{"slug = ?": ctx.Param("slug")}}
	result, err := h.serviceContainer.GetCRUDServices().SquadServices.GetList(c, squadQuery)
	if err != nil || len(result.List) == 0 {
		ctx.Redirect(http.StatusTemporaryRedirect, h.baseURL)
		return nil
	}
	squad := result.List[0]

	if collectionSlug := ctx.Param("collectionSlug"); collectionSlug != "" {
		collQuery := &pagination.ListQuery{ShowAll: true, Filters: map[string]interface{}{"slug = ?": collectionSlug}}
		collResult, _ := h.serviceContainer.GetCRUDServices().CollectionServices.GetList(c, collQuery)
		if len(collResult.List) > 0 {
			collection = collResult.List[0]
		}
	}

	return h.render(ctx, "squad_show", map[string]interface{}{
		"squad":      squad,
		"collection": collection,
	})
}

func (h handlers) viewShareDiagram(ctx echo.Context) error {
	c := context.Background()

	url := ctx.Param("url")
	if url == "" {
		return ctx.String(http.StatusBadRequest, "url is required")
	}

	decodeString, err := base64.RawURLEncoding.DecodeString(url)
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
