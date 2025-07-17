package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h handlers) routesViews() {
	h.server.GET("/squad/:slug", h.viewShowSquad)
	h.server.GET("/squad/:slug/collection/:collectionSlug", h.viewShowCollection)

	h.server.Static("v2/assets", "resources/views/vue/dist/assets")
	h.server.File("v2", "resources/views/vue/dist/index.html")
	h.server.File("v2/*", "resources/views/vue/dist/index.html")
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
