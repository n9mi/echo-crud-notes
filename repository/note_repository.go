package repository

import (
	"github.com/naomigrain/echo-crud-notes/helper"
	"github.com/naomigrain/echo-crud-notes/model/domain"
	"gorm.io/gorm"
)

type NoteRepository interface {
	FindAll(tx *gorm.DB, page int, pageSize int) ([]domain.ScanNote, error)
	IsExistById(tx *gorm.DB, id int) bool
	FindById(tx *gorm.DB, id int) (domain.ScanNote, error)
	Save(tx *gorm.DB, note domain.Note) (domain.Note, error)
	Delete(tx *gorm.DB, id int) error
}

type noteRepositoryImpl struct {
}

func NewNoteRepositoryImpl() *noteRepositoryImpl {
	return &noteRepositoryImpl{}
}

func (r *noteRepositoryImpl) FindAll(tx *gorm.DB, page int, pageSize int) ([]domain.ScanNote, error) {
	var note []domain.ScanNote
	if page > 0 && pageSize > 0 {
		tx = tx.Scopes(helper.Paginate(page, pageSize))
	}
	if err := tx.Model(&domain.Note{}).
		Order("notes.id asc").
		Select("notes.id, notes.title, notes.body, categories.name as category").
		Joins("inner join categories on categories.id = notes.category_id").
		Scan(&note).Error; err != nil {
		return note, err
	}

	return note, nil
}

func (r *noteRepositoryImpl) IsExistById(tx *gorm.DB, id int) bool {
	var count int64
	if tx.Model(&domain.Note{}).Where("id = ?", id).Count(&count); count == 0 {
		return false
	}

	return true
}

func (r *noteRepositoryImpl) FindById(tx *gorm.DB, id int) (domain.ScanNote, error) {
	var note domain.ScanNote
	if err := tx.Model(&domain.Note{}).
		Select("notes.id, notes.title, notes.body, categories.name as category").
		Where("notes.id = ?", id).
		Joins("inner join categories on categories.id = notes.category_id").
		Scan(&note).Error; err != nil {
		return note, err
	}

	return note, nil
}

func (r *noteRepositoryImpl) Save(tx *gorm.DB, note domain.Note) (domain.Note, error) {
	if err := tx.Save(&note).Error; err != nil {
		return note, err
	}

	return note, nil
}

func (r *noteRepositoryImpl) Delete(tx *gorm.DB, id int) error {
	if err := tx.Where("id = ?", id).Delete(&domain.Note{}).Error; err != nil {
		return err
	}

	return nil
}
