package repository

import (
	"github.com/naomigrain/echo-crud-notes/helper"
	"github.com/naomigrain/echo-crud-notes/model/domain"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(tx *gorm.DB, page int, pageSize int) ([]domain.Category, error)
	IsExistById(tx *gorm.DB, id int) bool
	FindById(tx *gorm.DB, id int) (domain.Category, error)
	Save(tx *gorm.DB, category domain.Category) (domain.Category, error)
	Delete(tx *gorm.DB, id int) error
}

type categoryRepositoryImpl struct {
}

func NewCategoryRepository() *categoryRepositoryImpl {
	return &categoryRepositoryImpl{}
}

func (r *categoryRepositoryImpl) FindAll(tx *gorm.DB, page int, pageSize int) ([]domain.Category, error) {
	var categories []domain.Category
	if page > 0 && pageSize > 0 {
		if err := tx.Scopes(helper.Paginate(page, pageSize)).Find(&categories).Error; err != nil {
			return categories, err
		}
	} else {
		if err := tx.Find(&categories).Error; err != nil {
			return categories, err
		}
	}

	return categories, nil
}

func (r *categoryRepositoryImpl) IsExistById(tx *gorm.DB, id int) bool {
	if result := tx.Model(&domain.Category{}).Select("id", id).RowsAffected; result == 0 {
		return false
	}

	return true
}

func (r *categoryRepositoryImpl) FindById(tx *gorm.DB, id int) (domain.Category, error) {
	var category domain.Category
	if err := tx.First(&category, id).Error; err != nil {
		return category, err
	}

	return category, nil
}

func (r *categoryRepositoryImpl) Save(tx *gorm.DB, category domain.Category) (domain.Category, error) {
	if err := tx.Save(&category).Error; err != nil {
		return category, err
	}

	return category, nil
}

func (r *categoryRepositoryImpl) Delete(tx *gorm.DB, id int) error {
	if err := tx.Delete(&domain.Category{}, id).Error; err != nil {
		return err
	}

	return nil
}
