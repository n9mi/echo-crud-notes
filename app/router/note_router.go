package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/naomigrain/echo-crud-notes/controller"
	"github.com/naomigrain/echo-crud-notes/repository"
	"github.com/naomigrain/echo-crud-notes/service"
	"gorm.io/gorm"
)

func NoteRouter(e *echo.Echo, mainUrl string, db *gorm.DB, validate *validator.Validate) {
	categoryRepository := repository.NewCategoryRepository()
	noteRepository := repository.NewNoteRepositoryImpl()
	service := service.NewNoteRepositoryImpl(db, validate, noteRepository, categoryRepository)
	controller := controller.NewNoteController(service)

	g := e.Group(mainUrl + "/notes")
	g.GET("", controller.GetAll)
	g.GET("/:id", controller.GetById)
	g.POST("", controller.Create)
	g.PUT("/:id", controller.Update)
	g.DELETE("/:id", controller.Delete)
}
