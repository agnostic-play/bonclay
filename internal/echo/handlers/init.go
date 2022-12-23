package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.linkaja.com/be/ditto/internal/services"
	"net/http"
)

type handlers struct {
	container services.ServiceContainer
}

func SetupHandlers(server *echo.Echo, container services.ServiceContainer) {
	handler := handlers{
		container: container,
	}
	server.GET("/index", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"name": "Bro!!",
		})
	}).Name = "foobar"
	server.GET("/readiness", handler.readinessHandlers)
}

func (h handlers) readinessHandlers(ctx echo.Context) (err error) {

	return ctx.JSON(200, h.container.Readiness())
}
