package handler

import (
	`net/http`
	`strconv`
	`strings`

	`github.com/agnostic-play/ditoo/internal/common/errs`
	`github.com/agnostic-play/ditoo/pkg/validator`
	`github.com/labstack/echo/v4`
)

type responseRender struct {
	BaseURL string
	Error   *string
	Data    interface{}
}

type responseBody struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type validateID struct {
	ID string `validate:"required,uuid"`
}

func (h handlers) render(ctx echo.Context, page string, data any) error {
	return ctx.Render(http.StatusOK, page, responseRender{
		BaseURL: h.baseURL,
		Error:   nil,
		Data:    data,
	})
}

func respJSON(ctx echo.Context, status int, data interface{}) error {
	return ctx.JSON(status, map[string]interface{}{
		"status":  status,
		"message": "success",
		"data":    data,
	})
}

func respErrJSON(ctx echo.Context, status int, err error) error {
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

func respErr(ctx echo.Context, status int, err error) error {
	msg := err.Error()

	if strings.Contains(err.Error(), "not found") {
		msg = "data not found"
		status = http.StatusNotFound

	}

	if getError, ok := err.(validator.BadRequestErrors); ok {
		return ctx.JSON(status, responseBody{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: getError.Error(),
			Data:    getError.GetBadRequestPayload(),
		})
	}

	if getError, ok := err.(errs.CustomError); ok {
		return ctx.JSON(status, responseBody{
			Status:  strconv.Itoa(getError.HttpCode),
			Message: getError.Error(),
			Data:    struct{}{},
		})
	}

	return ctx.JSON(status, map[string]interface{}{
		"status":  status,
		"message": msg,
		"data":    struct{}{},
	})
}

func validateRequest(ctx echo.Context, payload interface{}) (err error) {
	if err = ctx.Bind(&payload); err != nil {
		return
	}
	if err = ctx.Validate(payload); err != nil {
		return
	}
	return nil
}

func validateUUID(ctx echo.Context) (string, error) {

	var uuid validateID

	uuid.ID = ctx.Param("id")

	if err := ctx.Validate(uuid); err != nil {
		return "", err
	}

	return uuid.ID, nil
}
