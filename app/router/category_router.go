package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/naomigrain/echo-crud-notes/controller"
	"github.com/naomigrain/echo-crud-notes/repository"
	"github.com/naomigrain/echo-crud-notes/service"
	"gorm.io/gorm"
)

func CategoryRouter(e *echo.Echo, mainUrl string, db *gorm.DB, validate *validator.Validate) {
	repository := repository.NewCategoryRepository()
	service := service.NewCategoryService(db, validate, repository)
	controller := controller.NewCategoryController(service)

	g := e.Group(mainUrl + "/categories")
	g.GET("", controller.GetAll)
	g.GET("/:id", controller.GetById)
	g.POST("", controller.Create)
	g.PUT("/:id", controller.Update)
	g.DELETE("/:id", controller.Delete)
}
