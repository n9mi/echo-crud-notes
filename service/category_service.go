package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/naomigrain/echo-crud-notes/exception"
	"github.com/naomigrain/echo-crud-notes/model/domain"
	"github.com/naomigrain/echo-crud-notes/model/web"
	"github.com/naomigrain/echo-crud-notes/repository"
	"gorm.io/gorm"
)

type CategoryService interface {
	Create(category web.CategoryJSON) (web.CategoryJSON, error)
	GetById(id int) (web.CategoryJSON, error)
	GetAll(page int, pageSize int) ([]web.CategoryJSON, error)
	Update(category web.CategoryJSON) (web.CategoryJSON, error)
	Delete(id int) error
}

type categoryServiceImpl struct {
	DB         *gorm.DB
	Validate   *validator.Validate
	Repository repository.CategoryRepository
}

func NewCategoryService(db *gorm.DB, validate *validator.Validate, repository repository.CategoryRepository) *categoryServiceImpl {
	return &categoryServiceImpl{
		DB:         db,
		Validate:   validate,
		Repository: repository,
	}
}

func (s *categoryServiceImpl) GetAll(page int, pageSize int) ([]web.CategoryJSON, error) {
	var categories []web.CategoryJSON

	tx := s.DB.Begin()
	categoriesDom, errFind := s.Repository.FindAll(tx, page, pageSize)
	if errFind != nil {
		return categories, errFind
	}

	for _, cDom := range categoriesDom {
		categories = append(categories, web.CategoryJSON{
			ID:   cDom.ID,
			Name: cDom.Name,
		})
	}

	return categories, nil
}

func (s *categoryServiceImpl) GetById(id int) (web.CategoryJSON, error) {
	var category web.CategoryJSON

	tx := s.DB.Begin()
	categoryDom, errFind := s.Repository.FindById(tx, id)
	if errFind != nil {
		return category, &exception.NotFoundError{Entity: "category"}
	}

	category.ID = id
	category.Name = categoryDom.Name
	return category, nil
}

func (s *categoryServiceImpl) Create(category web.CategoryJSON) (web.CategoryJSON, error) {
	errValidate := s.Validate.Struct(category)
	if errValidate != nil {
		return category, errValidate
	}

	tx := s.DB.Begin()
	categoryDom, errCreate := s.Repository.Save(tx, domain.Category{
		Name: category.Name,
	})
	if errCreate != nil {
		errRollback := tx.Rollback().Error
		if errRollback != nil {
			return category, errRollback
		}
		return category, errCreate
	}

	if errCommit := tx.Commit().Error; errCommit != nil {
		return category, errCommit
	}

	category.ID = categoryDom.ID
	return category, nil
}

func (s *categoryServiceImpl) Update(category web.CategoryJSON) (web.CategoryJSON, error) {
	errValidate := s.Validate.Struct(category)
	if errValidate != nil {
		return category, errValidate
	}

	tx := s.DB.Begin()
	_, errFind := s.Repository.FindById(tx, category.ID)
	if errFind != nil {
		return category, &exception.NotFoundError{Entity: "category"}
	}

	_, errUpdate := s.Repository.Save(tx, domain.Category{
		ID:   category.ID,
		Name: category.Name,
	})
	if errUpdate != nil {
		errRollback := tx.Rollback().Error
		if errRollback != nil {
			return category, errRollback
		}
		return category, errUpdate
	}

	if errCommit := tx.Commit().Error; errCommit != nil {
		return category, errCommit
	}

	return category, nil
}

func (s *categoryServiceImpl) Delete(id int) error {
	tx := s.DB.Begin()
	_, errFind := s.Repository.FindById(tx, id)
	if errFind != nil {
		return &exception.NotFoundError{Entity: "category"}
	}

	errDel := s.Repository.Delete(tx, id)
	if errDel != nil {
		errRollback := tx.Rollback().Error
		if errRollback != nil {
			return errRollback
		}
		return errDel
	}

	if errCommit := tx.Commit().Error; errCommit != nil {
		return errCommit
	}

	return nil
}
