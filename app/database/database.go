package database

import (
	"fmt"
	"math/rand"

	"github.com/naomigrain/echo-crud-notes/config"
	"github.com/naomigrain/echo-crud-notes/helper"
	"github.com/naomigrain/echo-crud-notes/model/domain"
	"github.com/naomigrain/echo-crud-notes/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartConnection(dbConfig config.DBConfig) (*gorm.DB, error) {
	dsn := `host=` + dbConfig.DBHost + ` user=` + dbConfig.DBUsername +
		` password=` + dbConfig.DBPassword + ` dbname=` + dbConfig.DBName +
		` port=` + dbConfig.DBPort + ` sslmode=disable TimeZone=Asia/Jakarta`

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}

func Migrate(db *gorm.DB) {
	db.Migrator().CreateTable(&domain.Category{})
	db.Migrator().CreateTable(&domain.Note{})
}

func DropAll(db *gorm.DB) {
	db.Migrator().DropTable(&domain.Note{})
	db.Migrator().DropTable(&domain.Category{})
}

func CategorySeeder(db *gorm.DB, numRecords int) []domain.Category {
	categoryRepository := repository.NewCategoryRepository()
	var categoryList []domain.Category
	for i := 0; i < numRecords; i++ {
		categoryList = append(categoryList, domain.Category{
			Name: "Category " + helper.RandomString(rand.Intn(80)),
		})
	}
	for i, c := range categoryList {
		tx := db.Begin()
		categoryDom, _ := categoryRepository.Save(tx, c)
		tx.Commit()
		categoryList[i].ID = categoryDom.ID
	}
	return categoryList
}

func NoteSeeder(db *gorm.DB, categoryList []domain.Category, numRecords int) []domain.Note {
	noteRepository := repository.NewNoteRepositoryImpl()
	var noteList []domain.Note

	for i := 0; i < numRecords; i++ {
		tx := db.Begin()
		noteDom, _ := noteRepository.Save(tx, domain.Note{
			Title:      "Category " + helper.RandomString(rand.Intn(80)),
			Body:       "Body " + helper.RandomString(rand.Intn(100)),
			CategoryID: categoryList[rand.Intn(len(categoryList))].ID,
		})
		noteList = append(noteList, noteDom)
		tx.Commit()
	}

	return noteList
}

func DeleteCategoryRecords(db *gorm.DB) {
	db.Where("1=1").Delete(&domain.Category{})
}

func DeleteNoteRecords(db *gorm.DB) {
	db.Where("1=1").Delete(&domain.Note{})
}

func DeleteAllRecords(db *gorm.DB) {
	fmt.Println("Delete all records...")
	DeleteNoteRecords(db)
	DeleteCategoryRecords(db)
}
