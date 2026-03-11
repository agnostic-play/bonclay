package handler

import (
	`fmt`

	"github.com/agnostic-play/ditoo/internal/services"
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
	h.server.GET("/", func(ctx echo.Context) error {
		return h.render(ctx, "squadIndex", nil)
	})
	h.server.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, "hello :)")
	}).Name = "ping"

	crudServices := h.serviceContainer.GetCRUDServices()

	diagramHandlers := NewBaseCRUDHandlers(crudServices.DiagramServices)
	endpointHandlers := NewBaseCRUDHandlers(crudServices.EndpointServices)
	squadHandlers := NewBaseCRUDHandlers(crudServices.SquadServices)
	diagramCollHandlers := NewDiagramCollectionRoutes(crudServices.DiagramCollectionServices, crudServices.DiagramServices)
	collectionHandlers := NewCollectionRoutes(crudServices.CollectionServices, crudServices.EndpointServices, h.serviceContainer)
	scenarioHandler := NewScenarioRoutes(crudServices.ScenarioServices, h.serviceContainer)

	// Legacy routes (not yet migrated to v2)
	h.server.GET("/api/collection/custom_variable/list/:collectionSlug", h.getCustomVariable)
	h.server.POST("/api/collection/custom_variable/create", h.actCreateCustomVariable)
	h.server.GET("/api/squad/detail/:slug", h.getSquadDetail)
	h.server.Any("/mock/:collection/:path", h.actMock)

	h.server.GET("/share/diagram/:url", h.viewShareDiagram)
	h.server.GET("/squad/:slug", h.viewShowSquad)
	h.server.GET("/squad/:slug/collection/:collectionSlug", h.viewShowCollection)

	//h.server.Static("/v2/assets", "resources/views/spa/dist/assets")
	//h.server.File("/v2", "resources/views/spa/dist/index.html")
	//h.server.File("/v2/*", "resources/views/spa/dist/index.html")

	h.server.Static("/v2/assets", "resources/views/bonclay/public/bonclay/assets")
	h.server.File("/v2", "resources/views/bonclay/public/bonclay/index.html")
	h.server.File("/v2/*", "resources/views/bonclay/public/bonclay/index.html")

	v2 := h.server.Group("/v3")

	diagramHandlers.RegisterRoutes(v2)
	endpointHandlers.RegisterRoutes(v2)
	endpointHandlers.RegisterRoutes(v2)
	squadHandlers.RegisterRoutes(v2)
	diagramCollHandlers.RegisterRoutes(v2)
	collectionHandlers.RegisterRoutes(v2)
	scenarioHandler.RegisterRoutes(v2)

	fmt.Println("endpoints ===")
	for _, r := range h.server.Routes() {
		fmt.Println(r.Method, r.Path, r.Name)
	}
}
