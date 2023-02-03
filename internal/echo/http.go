package echo

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gitlab.linkaja.com/be/ditto/internal/config"
	handler "gitlab.linkaja.com/be/ditto/internal/echo/handlers"
	httpUtils "gitlab.linkaja.com/be/ditto/internal/echo/utils"
	"gitlab.linkaja.com/be/ditto/internal/services"
	"gitlab.linkaja.com/be/ditto/pkg/validator"
	"net/http"
)

func RunHttpServer(config *config.Config, container services.ServiceContainer) chan bool {
	e := echo.New()
	e.HideBanner = true
	e.DisableHTTP2 = true
	e.Validator = validator.NewValidator()
	e.Renderer = httpUtils.NewTemplates(config.App.TemplatePath)

	e.Static("/static", config.App.StaticPath)

	SetupMiddlewares(e)
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
