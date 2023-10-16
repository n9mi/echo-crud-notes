package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/naomigrain/echo-crud-notes/exception"
	"github.com/naomigrain/echo-crud-notes/model/domain"
	"github.com/naomigrain/echo-crud-notes/model/web"
	"github.com/naomigrain/echo-crud-notes/repository"
	"gorm.io/gorm"
)

type NoteService interface {
	GetAll(page int, pageSize int) ([]web.NoteResponse, error)
	GetById(id int) (web.NoteResponse, error)
	Create(note web.NoteRequest) (web.NoteResponse, error)
	Update(note web.NoteRequest) (web.NoteResponse, error)
	Delete(id int) error
}

type noteServiceImpl struct {
	DB                 *gorm.DB
	Validate           *validator.Validate
	NoteRepository     repository.NoteRepository
	CategoryRepository repository.CategoryRepository
}

func NewNoteRepositoryImpl(db *gorm.DB, validate *validator.Validate, noteRepository repository.NoteRepository,
	categoryRepository repository.CategoryRepository) *noteServiceImpl {
	return &noteServiceImpl{
		DB:                 db,
		Validate:           validate,
		NoteRepository:     noteRepository,
		CategoryRepository: categoryRepository,
	}
}

func (s *noteServiceImpl) GetAll(page int, pageSize int) ([]web.NoteResponse, error) {
	var notes []web.NoteResponse

	tx := s.DB.Begin()
	notesScan, errFind := s.NoteRepository.FindAll(tx, page, pageSize)
	if errFind != nil {
		return notes, errFind
	}

	for _, nS := range notesScan {
		notes = append(notes, web.NoteResponse{
			ID:       nS.ID,
			Title:    nS.Title,
			Body:     nS.Body,
			Category: nS.Category,
		})
	}

	return notes, nil
}

func (s *noteServiceImpl) GetById(id int) (web.NoteResponse, error) {
	var note web.NoteResponse

	tx := s.DB.Begin()
	noteScan, errFind := s.NoteRepository.FindById(tx, id)
	if errFind != nil || noteScan.Title == "" {
		return note, &exception.NotFoundError{Entity: "note"}
	}

	note = web.NoteResponse{
		ID:       id,
		Title:    noteScan.Title,
		Body:     noteScan.Body,
		Category: noteScan.Category,
	}
	return note, nil
}

func (s *noteServiceImpl) Create(note web.NoteRequest) (web.NoteResponse, error) {
	var noteResponse web.NoteResponse
	if errValidate := s.Validate.Struct(note); errValidate != nil {
		return noteResponse, errValidate
	}

	tx := s.DB.Begin()
	categoryDom, errFind := s.CategoryRepository.FindById(tx, note.CategoryId)
	if errFind != nil {
		return noteResponse, &exception.BadRequestError{Message: "category did not exists"}
	}

	noteDom, errSave := s.NoteRepository.Save(tx, domain.Note{
		Title:      note.Title,
		Body:       note.Body,
		CategoryID: note.CategoryId,
	})
	if errSave != nil {
		if errRollback := tx.Rollback().Error; errRollback != nil {
			return noteResponse, errRollback
		}
		return noteResponse, errSave
	}
	if errCommit := tx.Commit().Error; errCommit != nil {
		return noteResponse, errCommit
	}

	noteResponse = web.NoteResponse{
		ID:       noteDom.ID,
		Title:    note.Title,
		Body:     note.Body,
		Category: categoryDom.Name,
	}
	return noteResponse, nil
}

func (s *noteServiceImpl) Update(note web.NoteRequest) (web.NoteResponse, error) {
	var noteResponse web.NoteResponse
	if errVal := s.Validate.Struct(note); errVal != nil {
		return noteResponse, errVal
	}

	tx := s.DB.Begin()
	if isNoteExist := s.NoteRepository.IsExistById(tx, note.ID); !isNoteExist {
		return noteResponse, &exception.NotFoundError{Entity: "note"}
	}

	categoryDom, errFind := s.CategoryRepository.FindById(tx, note.CategoryId)
	if errFind != nil {
		return noteResponse, &exception.BadRequestError{Message: "category does not exists"}
	}

	noteDom, errUpdate := s.NoteRepository.Save(tx, domain.Note{
		ID:         note.ID,
		Title:      note.Title,
		Body:       note.Body,
		CategoryID: note.CategoryId,
	})
	if errUpdate != nil {
		if errRollback := tx.Rollback().Error; errRollback != nil {
			return noteResponse, errRollback
		}
		return noteResponse, errUpdate
	}
	if errCommit := tx.Commit().Error; errCommit != nil {
		return noteResponse, errCommit
	}

	noteResponse = web.NoteResponse{
		ID:       noteDom.ID,
		Title:    note.Title,
		Body:     note.Body,
		Category: categoryDom.Name,
	}
	return noteResponse, nil
}

func (s *noteServiceImpl) Delete(id int) error {
	tx := s.DB.Begin()
	if isNoteExist := s.NoteRepository.IsExistById(tx, id); !isNoteExist {
		return &exception.NotFoundError{Entity: "note"}
	}

	errDel := s.NoteRepository.Delete(tx, id)
	if errDel != nil {
		if errRollback := tx.Rollback().Error; errRollback != nil {
			return errRollback
		}
		return errDel
	}
	if errCommit := tx.Commit().Error; errCommit != nil {
		return errCommit
	}

	return nil
}
