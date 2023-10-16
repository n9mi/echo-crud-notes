package exception

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/naomigrain/echo-crud-notes/model/web"
)

func CustomErrorHandler(err error, c echo.Context) {
	var res web.ErrorResponse

	if _, ok := err.(*NotFoundError); ok {
		res.Code = http.StatusNotFound
		res.Status = "NOT FOUND"
		res.Message = err.Error()
	} else if errors.Is(err, echo.ErrNotFound) {
		res.Code = http.StatusNotFound
		res.Status = "NOT FOUND"
		res.Message = "Page does not exists"
	} else if _, ok := err.(*BadRequestError); ok {
		res.Code = http.StatusBadRequest
		res.Status = "BAD REQUEST"
		res.Message = err.Error()
	} else if castedErr, ok := err.(validator.ValidationErrors); ok {
		res.Code = http.StatusBadRequest
		res.Status = "BAD REQUEST"
		for _, e := range castedErr {
			switch e.Tag() {
			case "required":
				res.Message = fmt.Sprintf("%s is required", e.Field())
			case "max":
				res.Message = fmt.Sprintf("%s is should below than %s characters", e.Field(), e.Param())
			case "min":
				res.Message = fmt.Sprintf("%s is should more than %s characters", e.Field(), e.Param())
			case "gte":
				res.Message = fmt.Sprintf("%s is should greater than %s", e.Field(), e.Param())
			}
		}
	} else {
		res.Code = http.StatusInternalServerError
		res.Status = "FAIL"
		res.Message = err.Error()
	}

	c.Logger().Error(err)
	c.JSON(res.Code, res)
}
