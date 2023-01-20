package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.linkaja.com/be/ditto/internal/services"
)

type HandlersInit interface {
	Routes()
}

type handlers struct {
	serviceContainer services.ServiceContainer
	server           *echo.Echo
	baseURL          string
}

type responseRender struct {
	BaseURL string
	Error   *string
	Data    interface{}
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

	h.routesSquad()
	h.routesCollection()

}
