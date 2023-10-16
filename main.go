package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/naomigrain/echo-crud-notes/app/database"
	"github.com/naomigrain/echo-crud-notes/app/router"
	"github.com/naomigrain/echo-crud-notes/config"
)

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func main() {
	dbConfig := config.GetDBConfig(true)
	db, errConn := database.StartConnection(dbConfig)
	if errConn != nil {
		panic(errConn)
	}
	database.DropAll(db)
	database.Migrate(db)
	categoryList := database.CategorySeeder(db, 3)
	database.NoteSeeder(db, categoryList, 5)

	validate := validator.New()

	e := router.InitializeEcho()
	router.AssignRouter(e, db, validate)

	port := config.GetAppConfig(true).AppPort
	e.Logger.Fatal(e.Start(":" + port))
}
