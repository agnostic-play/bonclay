package echo

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.linkaja.com/be/ditto/internal/constant"
	httpUtils "gitlab.linkaja.com/be/ditto/internal/echo/utils"
)

func SetupMiddlewares(server *echo.Echo) {
	server.Pre(middleware.RemoveTrailingSlash())
	server.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			session := httpUtils.NewSessionRequest()
			session.IP = ctx.RealIP()
			session.URL = ctx.Request().URL.String()
			session.Method = ctx.Request().Method

			ctx.Set(constant.AppSessionRequest, session)
			return h(ctx)
		}
	})

	//server.Use(middleware.Recover())
}
