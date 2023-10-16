package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/naomigrain/echo-crud-notes/exception"
	"github.com/naomigrain/echo-crud-notes/model/web"
	"github.com/naomigrain/echo-crud-notes/service"
)

type NoteController interface {
	GetAll(c echo.Context) error
	GetById(c echo.Context) error
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type noteControllerImpl struct {
	Service service.NoteService
}

func NewNoteController(service service.NoteService) *noteControllerImpl {
	return &noteControllerImpl{
		Service: service,
	}
}

func (ct *noteControllerImpl) GetAll(c echo.Context) error {
	page := c.QueryParam("page")
	pageSize := c.QueryParam("pageSize")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	noteRes, errFind := ct.Service.GetAll(pageInt, pageSizeInt)
	if errFind != nil {
		return errFind
	}

	response := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   noteRes,
	}
	return c.JSON(http.StatusOK, response)
}

func (ct *noteControllerImpl) GetById(c echo.Context) error {
	id := c.Param("id")
	idInt, errConv := strconv.Atoi(id)
	if errConv != nil {
		return &exception.NotFoundError{Entity: "note"}
	}

	noteRes, errFind := ct.Service.GetById(idInt)
	if errFind != nil {
		return &exception.NotFoundError{Entity: "note"}
	}

	response := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   noteRes,
	}
	return c.JSON(http.StatusOK, response)
}

func (ct *noteControllerImpl) Create(c echo.Context) error {
	noteReq := new(web.NoteRequest)
	if errBind := c.Bind(noteReq); errBind != nil {
		return &exception.BadRequestError{Message: errBind.Error()}
	}

	noteRes, errCreate := ct.Service.Create(*noteReq)
	if errCreate != nil {
		return errCreate
	}

	response := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   noteRes,
	}
	return c.JSON(http.StatusOK, response)
}

func (ct *noteControllerImpl) Update(c echo.Context) error {
	id := c.Param("id")
	idInt, errConv := strconv.Atoi(id)
	if errConv != nil {
		return errConv
	}

	noteReq := new(web.NoteRequest)
	if errBind := c.Bind(noteReq); errBind != nil {
		return &exception.BadRequestError{Message: errBind.Error()}
	}

	noteReq.ID = idInt
	noteRes, errUpdate := ct.Service.Update(*noteReq)
	if errUpdate != nil {
		return errUpdate
	}

	response := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   noteRes,
	}
	return c.JSON(http.StatusOK, response)
}

func (ct *noteControllerImpl) Delete(c echo.Context) error {
	id := c.Param("id")
	idInt, errConv := strconv.Atoi(id)
	if errConv != nil {
		return &exception.NotFoundError{Entity: "note"}
	}

	errDel := ct.Service.Delete(idInt)
	if errDel != nil {
		return errDel
	}

	response := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
	}
	return c.JSON(http.StatusOK, response)
}
