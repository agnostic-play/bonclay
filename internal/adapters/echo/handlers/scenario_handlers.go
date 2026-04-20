package handler

import (
	"context"
	"fmt"
	"net/http"

	"berlin.allobank.com/tools/bonclay/internal/entities"
	"berlin.allobank.com/tools/bonclay/internal/services"
	api_mock_services "berlin.allobank.com/tools/bonclay/internal/services/api-mock-services"
	crud_services "berlin.allobank.com/tools/bonclay/internal/services/crud-services"
	"github.com/labstack/echo/v4"
)

// ScenarioHandlers extends BaseCRUDHandlers with scenario-specific endpoints
type ScenarioHandlers interface {
	BaseCRUDHandlers[entities.ScenarioEntity]
}

type scenarioHandlers struct {
	baseCRUDHandlers BaseCRUDHandlers[entities.ScenarioEntity]
	serviceContainer services.ServiceContainer
}

// NewScenarioHandlers creates a new ScenarioHandlers instance
func NewScenarioRoutes(
	scenarioServices crud_services.BaseCRUDService[entities.ScenarioEntity],
	serviceContainer services.ServiceContainer,
) BaseCRUDHandlers[entities.ScenarioEntity] {
	return &scenarioHandlers{
		baseCRUDHandlers: NewBaseCRUDHandlers(scenarioServices),
		serviceContainer: serviceContainer,
	}
}

func (h scenarioHandlers) RegisterRoutes(e *echo.Group) {
	// Register standard CRUD routes
	h.baseCRUDHandlers.RegisterRoutes(e)

	// Register custom routes
	path := h.baseCRUDHandlers.GetPath()
	e.POST(fmt.Sprintf("/%s/set_active", path), h.setActiveScenario)
	e.POST(fmt.Sprintf("/%s/:id/remove_active", path), h.removeActiveScenario)
}

func (h scenarioHandlers) GetPath() string {
	return h.baseCRUDHandlers.GetPath()
}

func (h scenarioHandlers) GetServices() crud_services.BaseCRUDService[entities.ScenarioEntity] {
	return h.baseCRUDHandlers.GetServices()
}

// setActiveScenario sets a scenario as active for its endpoint
func (h scenarioHandlers) setActiveScenario(ctx echo.Context) error {
	var req api_mock_services.SetActiveScenarioReq

	if err := validateRequest(ctx, &req); err != nil {
		return respErr(ctx, http.StatusBadRequest, err)
	}

	err := h.serviceContainer.GetEndpointScenarioService().SetActiveResponse(context.Background(), req)
	if err != nil {
		return respErr(ctx, http.StatusInternalServerError, err)
	}

	return respJSON(ctx, http.StatusOK, req)
}

// removeActiveScenario removes the active scenario from an endpoint
func (h scenarioHandlers) removeActiveScenario(ctx echo.Context) error {
	var req api_mock_services.SetActiveScenarioReq
	req.EndpointID = ctx.Param("id")

	err := h.serviceContainer.GetEndpointScenarioService().RemoveScenario(context.Background(), req)
	if err != nil {
		return respErr(ctx, http.StatusInternalServerError, err)
	}

	return respJSON(ctx, http.StatusOK, req)
}
