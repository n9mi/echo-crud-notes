package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/naomigrain/echo-crud-notes/exception"
	"github.com/naomigrain/echo-crud-notes/model/web"
	"github.com/naomigrain/echo-crud-notes/service"
)

type CategoryController interface {
	Create(c echo.Context) error
	GetById(c echo.Context) error
	GetAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type categoryControllerImpl struct {
	Service service.CategoryService
}

func NewCategoryController(service service.CategoryService) *categoryControllerImpl {
	return &categoryControllerImpl{
		Service: service,
	}
}

func (ct *categoryControllerImpl) GetAll(c echo.Context) error {
	page := c.QueryParam("page")
	pageSize := c.QueryParam("pageSize")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	categories, errFind := ct.Service.GetAll(pageInt, pageSizeInt)
	if errFind != nil {
		return errFind
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   categories,
	}
	return c.JSON(http.StatusOK, res)
}

func (ct *categoryControllerImpl) GetById(c echo.Context) error {
	id := c.Param("id")
	idInt, errConv := strconv.Atoi(id)
	if errConv != nil {
		return &exception.NotFoundError{Entity: "category"}
	}

	categoryRes, errFind := ct.Service.GetById(idInt)
	if errFind != nil {
		return errFind
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   categoryRes,
	}
	return c.JSON(http.StatusInternalServerError, res)
}

func (ct *categoryControllerImpl) Create(c echo.Context) error {
	categoryReq := new(web.CategoryJSON)
	if errBind := c.Bind(categoryReq); errBind != nil {
		return &exception.BadRequestError{Message: errBind.Error()}
	}

	categoryRes, errCreate := ct.Service.Create(*categoryReq)
	if errCreate != nil {
		return errCreate
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   categoryRes,
	}
	return c.JSON(http.StatusOK, res)
}

func (ct *categoryControllerImpl) Update(c echo.Context) error {
	id := c.Param("id")
	idInt, errConv := strconv.Atoi(id)
	if errConv != nil {
		return &exception.NotFoundError{Entity: "category"}
	}

	categoryReq := new(web.CategoryJSON)
	if errBind := c.Bind(categoryReq); errBind != nil {
		return &exception.BadRequestError{Message: errBind.Error()}
	}

	categoryReq.ID = idInt
	categoryRes, errUpdate := ct.Service.Update(*categoryReq)
	if errUpdate != nil {
		return errUpdate
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   categoryRes,
	}
	return c.JSON(http.StatusOK, res)
}

func (ct *categoryControllerImpl) Delete(c echo.Context) error {
	id := c.Param("id")
	idInt, errConv := strconv.Atoi(id)
	if errConv != nil {
		return &exception.NotFoundError{Entity: "category"}
	}

	errDel := ct.Service.Delete(idInt)
	if errDel != nil {
		return errDel
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
	}
	return c.JSON(http.StatusOK, res)
}
