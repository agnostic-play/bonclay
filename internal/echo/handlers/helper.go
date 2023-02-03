package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.linkaja.com/be/ditto/internal/errs"
	"net/http"
	"strings"
)

type validateID struct {
	ID string `validate:"required,uuid"`
}

func (h handlers) render(ctx echo.Context, page string, data interface{}) error {
	return ctx.Render(http.StatusOK, page, responseRender{
		BaseURL: h.baseURL,
		Error:   nil,
		Data:    data,
	})
}

func (h handlers) json(ctx echo.Context, status int, data interface{}) error {
	return ctx.JSON(status, map[string]interface{}{
		"status":  status,
		"message": "success",
		"data":    data,
	})
}

func (h handlers) errorJson(ctx echo.Context, status int, err error) error {
	msg := err.Error()

	if strings.Contains(err.Error(), "not found") {
		msg = "data not found"
		status = http.StatusNotFound

	}

	if getError, ok := err.(errs.CustomError); ok {
		status = getError.HttpCode
		msg = getError.Error()
		err = nil
	}

	return ctx.JSON(status, map[string]interface{}{
		"status":  status,
		"message": msg,
		"data":    err,
	})
}

func (h handlers) validateRequest(ctx echo.Context, payload interface{}) (err error) {
	if err = ctx.Bind(&payload); err != nil {
		return
	}
	if err = ctx.Validate(payload); err != nil {
		return
	}
	return nil
}

func (h handlers) validateUUID(ctx echo.Context) (string, error) {

	var uuid validateID

	uuid.ID = ctx.Param("id")

	if err := ctx.Validate(uuid); err != nil {
		return "", err
	}

	return uuid.ID, nil
}
