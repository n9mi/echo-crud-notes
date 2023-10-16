package test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/naomigrain/echo-crud-notes/app/database"
	"github.com/naomigrain/echo-crud-notes/app/router"
	"github.com/naomigrain/echo-crud-notes/config"
	"gorm.io/gorm"
)

var e *echo.Echo
var db *gorm.DB
var errConn error

func init() {
	dbConfig := config.GetDBConfig(true)
	db, errConn = database.StartConnection(dbConfig)
	if errConn != nil {
		panic(errConn)
	}
	database.DropAll(db)
	database.Migrate(db)
	database.CategorySeeder(db, 5)

	validate := validator.New()

	e = router.InitializeEcho()
	router.AssignRouter(e, db, validate)
}

func newTestRequest(url string, method string, requestBody string) *http.Request {
	request := httptest.NewRequest(method, url, strings.NewReader(requestBody))
	request.Header.Add("Content-Type", "application/json")

	return request
}
