package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"gitlab.linkaja.com/be/ditto/internal/services"
	"log"
	"net/http"
)

func (h handlers) routesSquad() {
	h.server.POST("/api/squad/create", h.actCreateSquad)
	h.server.DELETE("/api/squad/delete/:id", h.actDeleteSquad)
	h.server.GET("/api/squad", h.getListSquad)

	h.server.GET("/squad/:slug", h.viewShowSquad)
	h.server.GET("/squad/:slug/collection/:collectionSlug", h.viewShowSquad)
}

func (h handlers) viewShowSquad(ctx echo.Context) error {
	var collection interface{}
	c := context.Background()

	resp, err := h.serviceContainer.GetSquad(c, ctx.Param("slug"))
	if err != nil {
		return h.render(ctx, "squad_show", nil)
	}

	if collectionSlug := ctx.Param("collectionSlug"); collectionSlug != "" {
		collection, _ = h.serviceContainer.GetCollection(c, collectionSlug)
	}
	log.Println(collection)
	log.Println(ctx.Param("collectionSlug"))

	return h.render(ctx, "squad_show", map[string]interface{}{
		"squad":      resp,
		"collection": collection,
	})
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
