package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/naomigrain/echo-crud-notes/exception"
	"gorm.io/gorm"
)

func InitializeEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	e.HTTPErrorHandler = exception.CustomErrorHandler

	return e
}

func AssignRouter(e *echo.Echo, db *gorm.DB, validate *validator.Validate) {
	mainUrl := "/api"
	CategoryRouter(e, mainUrl, db, validate)
	NoteRouter(e, mainUrl, db, validate)
}
