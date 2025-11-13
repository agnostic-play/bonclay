package echo

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/agnostic-play/ditoo/internal/common/constant"
	"github.com/agnostic-play/ditoo/internal/config"
	handler "github.com/agnostic-play/ditoo/internal/echo/handlers"
	httpUtils "github.com/agnostic-play/ditoo/internal/echo/utils"
	"github.com/agnostic-play/ditoo/internal/services"
	"github.com/agnostic-play/ditoo/pkg/validator"
)

func RunHttpServer(config *config.Config, container services.ServiceContainer) chan bool {
	e := echo.New()
	e.HideBanner = true
	e.DisableHTTP2 = true
	e.Validator = validator.NewValidator()
	e.Renderer = httpUtils.NewTemplates(config.App.TemplatePath)

	e.Static("/static", config.App.StaticPath)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			session := httpUtils.NewSessionRequest()
			session.IP = ctx.RealIP()
			session.URL = ctx.Request().URL.String()
			session.Method = ctx.Request().Method

			ctx.Set(constant.AppSessionRequest, session)
			return h(ctx)
		}
	})
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*",
			"http://localhost:*",
			"http://localhost:5173/*",
			"https://api-mockup.alexandria.allobank.local/*",
			"http://api-mockup.alexandria.allobank.local/*",
		},
	}))

	handler.NewSetupHandlers(e, config.App.BaseURL, container).Routes()

	go func() {
		e.HideBanner = true
		address := fmt.Sprintf("%v:%v", config.App.Host, config.App.Port)
		if err := e.Start(address); err != nil {

			// ErrServerClosed is expected behaviour when exiting app
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("%v server, %v", "server", err)
			}
			log.Infof("%v server, %v", "server", err)
		}
	}()

	exitSignal := make(chan bool)
	go func() {
		<-exitSignal
		log.Info("stopping http server")
		if err := e.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed stopping server, %v", err)
		}
		log.Info("finished stopping http server")
		exitSignal <- true // Send signal already finish the job
	}()

	return exitSignal
}
