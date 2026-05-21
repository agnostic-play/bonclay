package handler

import (
	"berlin.allobank.com/tools/bonclay/internal/services"
	"github.com/labstack/echo/v4"
)

type HandlersInit interface {
	Routes()
}

type handlers struct {
	serviceContainer services.ServiceContainer
	server           *echo.Echo
	baseURL          string
}

func NewSetupHandlers(server *echo.Echo, baseURL string, container services.ServiceContainer) HandlersInit {
	return &handlers{
		server:           server,
		serviceContainer: container,
		baseURL:          baseURL,
	}
}

func (h handlers) Routes() {

	h.server.Static("/assets", "resources/views/bonclay/public/bonclay/assets")
	h.server.Static("/static", "resources/public/assets")
	h.server.File("/", "resources/views/bonclay/public/bonclay/index.html")
	h.server.File("/*", "resources/views/bonclay/public/bonclay/index.html")
	h.server.GET("/old", func(ctx echo.Context) error {
		return h.render(ctx, "squadIndex", nil)
	})
	h.server.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, "hello :)")
	}).Name = "ping"

	crudServices := h.serviceContainer.GetCRUDServices()

	diagramHandlers := NewBaseCRUDHandlers(crudServices.DiagramServices)
	endpointHandlers := NewEndpointHandlers(crudServices.EndpointServices)
	squadHandlers := NewBaseCRUDHandlers(crudServices.SquadServices)
	diagramCollHandlers := NewDiagramCollectionHandlers(crudServices.DiagramCollectionServices, crudServices.DiagramServices)
	collectionHandler := NewCollectionHandlers(crudServices.CollectionServices, crudServices.EndpointServices, h.serviceContainer)
	scenarioHandler := NewScenarioRoutes(crudServices.ScenarioServices, h.serviceContainer)
	customVariableHandlers := NewBaseCRUDHandlers(crudServices.CustomVariableServices)
	kongHandler := NewKongHandlers(h.serviceContainer.GetKongService())

	// Legacy routes (not yet migrated to v2)
	h.server.GET("/api/collection/custom_variable/list/:collectionSlug", h.getCustomVariable)
	h.server.POST("/api/collection/custom_variable/create", h.actCreateCustomVariable)

	h.server.Any("/mock/:collection/:path", h.actMock)

	h.server.GET("/share/diagram/:url", h.viewShareDiagram)
	h.server.GET("/squad/:slug", h.viewShowSquad)
	h.server.GET("/squad/:slug/collection/:collectionSlug", h.viewShowCollection)

	h.server.Static("/v2/assets", "resources/views/spa/dist/assets")
	h.server.File("/v2", "resources/views/spa/dist/index.html")
	h.server.File("/v2/*", "resources/views/spa/dist/index.html")

	apiV2 := h.server.Group("/api/v2")
	diagramHandlers.RegisterRoutes(apiV2)
	endpointHandlers.RegisterRoutes(apiV2)
	squadHandlers.RegisterRoutes(apiV2)
	diagramCollHandlers.RegisterRoutes(apiV2)
	collectionHandler.RegisterRoutes(apiV2)
	scenarioHandler.RegisterRoutes(apiV2)
	customVariableHandlers.RegisterRoutes(apiV2)
	kongHandler.RegisterRoutes(apiV2)

	apiV2.GET("/squad/detail/:slug", h.getSquadDetail)

	apiV2.GET("/collections/:slug/history", h.getHistory)
	apiV2.DELETE("/collections/:slug/history", h.clearHistory)

	// Environment variables for a collection
	apiV2.GET("/collections/:slug/variables", h.listCollectionVariables)
	apiV2.POST("/collections/:slug/variables", h.createCollectionVariable)
	apiV2.PATCH("/collections/variables/:id", h.updateCollectionVariable)
	apiV2.DELETE("/collections/variables/:id", h.deleteCollectionVariable)

}
