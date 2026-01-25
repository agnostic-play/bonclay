package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/agnostic-play/ditoo/internal/entities"
	"github.com/agnostic-play/ditoo/internal/services"
	api_mock_services "github.com/agnostic-play/ditoo/internal/services/api-mock-services"
	crud_services "github.com/agnostic-play/ditoo/internal/services/crud-services"
	"github.com/labstack/echo/v4"
)

// ScenarioRoutes extends BaseCRUDRoutes with scenario-specific endpoints
type ScenarioRoutes interface {
	BaseCRUDRoutes[entities.ScenarioEntity]
}

type scenarioRoutes struct {
	baseCRUDRoutes   BaseCRUDRoutes[entities.ScenarioEntity]
	serviceContainer services.ServiceContainer
}

// NewScenarioRoutes creates a new ScenarioRoutes instance
func NewScenarioRoutes(
	scenarioServices crud_services.BaseCRUDService[entities.ScenarioEntity],
	serviceContainer services.ServiceContainer,
) BaseCRUDRoutes[entities.ScenarioEntity] {
	return &scenarioRoutes{
		baseCRUDRoutes:   NewBaseCRUDRoutes(scenarioServices),
		serviceContainer: serviceContainer,
	}
}

func (h scenarioRoutes) RegisterRoutes(e *echo.Group) {
	// Register standard CRUD routes
	h.baseCRUDRoutes.RegisterRoutes(e)

	// Register custom routes
	path := h.baseCRUDRoutes.GetPath()
	e.POST(fmt.Sprintf("/%s/set_active", path), h.setActiveScenario)
	e.POST(fmt.Sprintf("/%s/:id/remove_active", path), h.removeActiveScenario)
}

func (h scenarioRoutes) GetPath() string {
	return h.baseCRUDRoutes.GetPath()
}

func (h scenarioRoutes) GetServices() crud_services.BaseCRUDService[entities.ScenarioEntity] {
	return h.baseCRUDRoutes.GetServices()
}

// setActiveScenario sets a scenario as active for its endpoint
func (h scenarioRoutes) setActiveScenario(ctx echo.Context) error {
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
func (h scenarioRoutes) removeActiveScenario(ctx echo.Context) error {
	var req api_mock_services.SetActiveScenarioReq
	req.EndpointID = ctx.Param("id")

	err := h.serviceContainer.GetEndpointScenarioService().RemoveScenario(context.Background(), req)
	if err != nil {
		return respErr(ctx, http.StatusInternalServerError, err)
	}

	return respJSON(ctx, http.StatusOK, req)
}
