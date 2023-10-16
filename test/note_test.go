package test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/naomigrain/echo-crud-notes/app/database"
	"github.com/naomigrain/echo-crud-notes/helper"
	"github.com/naomigrain/echo-crud-notes/model/domain"
	"github.com/naomigrain/echo-crud-notes/model/web"
	"github.com/stretchr/testify/require"
)

type testNoteJSON struct {
	Code   int              `json:"code"`
	Status string           `json:"status"`
	Data   web.NoteResponse `json:"data"`
}

type testNoteListJSON struct {
	Code   int                `json:"code"`
	Status string             `json:"status"`
	Data   []web.NoteResponse `json:"data"`
}

var noteUrl string = "http://127.0.0.1:8000/api/notes"
var categoryList []domain.Category

func findCategoryInList(id int, categoryList []domain.Category) domain.Category {
	for _, c := range categoryList {
		if c.ID == id {
			return c
		}
	}

	return domain.Category{}
}

func TestGetNotes(t *testing.T) {
	categoryList = database.CategorySeeder(db, 3)
	noteList := database.NoteSeeder(db, categoryList, 5)
	defer database.DeleteAllRecords(db)

	t.Run("Category_Get_By_Id_Success", func(t *testing.T) {
		randNote := noteList[rand.Intn(len(noteList))]
		category := findCategoryInList(randNote.CategoryID, categoryList)
		getByIdUrl := noteUrl + "/" + strconv.Itoa(randNote.ID)

		request := newTestRequest(getByIdUrl, http.MethodGet, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Body

		responseBody, _ := io.ReadAll(response)
		var noteResponse testNoteJSON
		json.Unmarshal(responseBody, &noteResponse)

		require.Equal(t, http.StatusOK, noteResponse.Code)
		require.Equal(t, "OK", noteResponse.Status)
		require.Equal(t, randNote.ID, noteResponse.Data.ID)
		require.Equal(t, randNote.Title, noteResponse.Data.Title)
		require.Equal(t, randNote.Body, noteResponse.Data.Body)
		require.Equal(t, category.Name, noteResponse.Data.Category)
	})

	t.Run("Category_Get_By_Id_NotFound_Fail", func(t *testing.T) {
		getByIdUrl := noteUrl + "/" + strconv.Itoa(9999999)

		request := newTestRequest(getByIdUrl, http.MethodGet, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Body

		responseBody, _ := io.ReadAll(response)
		var noteResponse web.ErrorResponse
		json.Unmarshal(responseBody, &noteResponse)

		require.Equal(t, http.StatusNotFound, noteResponse.Code)
		require.Equal(t, "NOT FOUND", noteResponse.Status)
		require.Equal(t, "note not found", noteResponse.Message)
	})

	t.Run("Category_Get_All_Success", func(t *testing.T) {
		request := newTestRequest(noteUrl, http.MethodGet, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Body

		responseBody, _ := io.ReadAll(response)
		var noteResponse testNoteListJSON
		json.Unmarshal(responseBody, &noteResponse)

		require.Equal(t, http.StatusOK, noteResponse.Code)
		require.Equal(t, "OK", noteResponse.Status)
		require.Equal(t, len(noteList), len(noteResponse.Data))
		for i, nD := range noteResponse.Data {
			require.Equal(t, noteList[i].ID, nD.ID)
			require.Equal(t, noteList[i].Title, nD.Title)
			require.Equal(t, noteList[i].Body, nD.Body)
			categoryNote := findCategoryInList(noteList[i].CategoryID, categoryList)
			require.Equal(t, categoryNote.Name, nD.Category)
		}
	})
}

func TestCreateNote(t *testing.T) {
	createUrl := noteUrl
	categoryList = database.CategorySeeder(db, 3)
	defer database.DeleteAllRecords(db)

	t.Run("Note_Create_Success", func(t *testing.T) {
		randString := helper.RandomString(rand.Intn(80))
		noteTitle := "Title " + randString
		noteBody := "Body " + randString
		randCatId := rand.Intn(len(categoryList))
		noteCategory := categoryList[randCatId]

		requestBody := fmt.Sprintf(
			`{"title": "%s", "body": "%s", "id_category": %d}`,
			noteTitle, noteBody, noteCategory.ID)
		request := newTestRequest(createUrl, http.MethodPost, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate testNoteJSON
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusOK, responseCreate.Code)
		require.Equal(t, "OK", responseCreate.Status)
		require.Equal(t, noteTitle, responseCreate.Data.Title)
		require.Equal(t, noteBody, responseCreate.Data.Body)
		require.Equal(t, noteCategory.Name, responseCreate.Data.Category)
	})

	t.Run("Note_Create_BadRequest_Fail", func(t *testing.T) {
		randString := helper.RandomString(rand.Intn(80))
		noteTitle := "T"
		noteBody := "Body " + randString
		randCatId := rand.Intn(len(categoryList))
		noteCategory := categoryList[randCatId]

		requestBody := fmt.Sprintf(
			`{"title": "%s", "body": "%s", "id_category": %d}`,
			noteTitle, noteBody, noteCategory.ID)
		request := newTestRequest(createUrl, http.MethodPost, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate web.ErrorResponse
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusBadRequest, responseCreate.Code)
		require.Equal(t, "BAD REQUEST", responseCreate.Status)
		require.Equal(t, "Title is should more than 2 characters", responseCreate.Message)
	})

	t.Run("Note_Create_BadRequest2_Fail", func(t *testing.T) {
		randString := helper.RandomString(rand.Intn(80))
		noteTitle := "Title " + randString
		noteBody := "B"
		randCatId := rand.Intn(len(categoryList))
		noteCategory := categoryList[randCatId]

		requestBody := fmt.Sprintf(
			`{"title": "%s", "body": "%s", "id_category": %d}`,
			noteTitle, noteBody, noteCategory.ID)
		request := newTestRequest(createUrl, http.MethodPost, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate web.ErrorResponse
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusBadRequest, responseCreate.Code)
		require.Equal(t, "BAD REQUEST", responseCreate.Status)
		require.Equal(t, "Body is should more than 2 characters", responseCreate.Message)
	})

	t.Run("Note_Create_BadRequest3_Fail", func(t *testing.T) {
		randString := helper.RandomString(rand.Intn(80))
		noteTitle := "Title " + randString
		noteBody := "Body " + randString

		requestBody := fmt.Sprintf(
			`{"title": "%s", "body": "%s", "id_category": %d}`,
			noteTitle, noteBody, 9999999)
		request := newTestRequest(createUrl, http.MethodPost, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate web.ErrorResponse
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusBadRequest, responseCreate.Code)
		require.Equal(t, "BAD REQUEST", responseCreate.Status)
		require.Equal(t, "category did not exists", responseCreate.Message)
	})
}

func TestUpdateNote(t *testing.T) {
	defer database.DeleteAllRecords(db)
	categoryList = database.CategorySeeder(db, 3)
	noteList := database.NoteSeeder(db, categoryList, 5)

	t.Run("Note_Update_Success", func(t *testing.T) {
		note := noteList[rand.Intn(len(noteList))]
		updateUrl := noteUrl + "/" + strconv.Itoa(note.ID)

		randString := helper.RandomString(rand.Intn(80))
		noteTitle := "Title " + randString
		noteBody := "Body " + randString
		randCatId := rand.Intn(len(categoryList))
		noteCategory := categoryList[randCatId]

		requestBody := fmt.Sprintf(
			`{"title": "%s", "body": "%s", "id_category": %d}`,
			noteTitle, noteBody, noteCategory.ID)
		request := newTestRequest(updateUrl, http.MethodPut, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate testNoteJSON
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusOK, responseCreate.Code)
		require.Equal(t, "OK", responseCreate.Status)
		require.Equal(t, note.ID, responseCreate.Data.ID)
		require.Equal(t, noteTitle, responseCreate.Data.Title)
		require.Equal(t, noteBody, responseCreate.Data.Body)
		require.Equal(t, noteCategory.Name, responseCreate.Data.Category)
	})

	t.Run("Note_Update_NotFound_Fail", func(t *testing.T) {
		updateUrl := noteUrl + "/" + strconv.Itoa(9999999)

		randString := helper.RandomString(rand.Intn(80))
		noteTitle := "Title " + randString
		noteBody := "Body " + randString
		randCatId := rand.Intn(len(categoryList))
		noteCategory := categoryList[randCatId]

		requestBody := fmt.Sprintf(
			`{"title": "%s", "body": "%s", "id_category": %d}`,
			noteTitle, noteBody, noteCategory.ID)
		request := newTestRequest(updateUrl, http.MethodPut, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate web.ErrorResponse
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusNotFound, responseCreate.Code)
		require.Equal(t, "NOT FOUND", responseCreate.Status)
		require.Equal(t, "note not found", responseCreate.Message)
	})

	t.Run("Note_Update_BadRequest_Fail", func(t *testing.T) {
		note := noteList[rand.Intn(len(noteList))]
		updateUrl := noteUrl + "/" + strconv.Itoa(note.ID)

		randString := helper.RandomString(rand.Intn(80))
		noteTitle := "T"
		noteBody := "Body " + randString
		randCatId := rand.Intn(len(categoryList))
		noteCategory := categoryList[randCatId]

		requestBody := fmt.Sprintf(
			`{"title": "%s", "body": "%s", "id_category": %d}`,
			noteTitle, noteBody, noteCategory.ID)
		request := newTestRequest(updateUrl, http.MethodPut, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate web.ErrorResponse
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusBadRequest, responseCreate.Code)
		require.Equal(t, "BAD REQUEST", responseCreate.Status)
		require.Equal(t, "Title is should more than 2 characters", responseCreate.Message)
	})

	t.Run("Note_Update_BadRequest2_Fail", func(t *testing.T) {
		note := noteList[rand.Intn(len(noteList))]
		updateUrl := noteUrl + "/" + strconv.Itoa(note.ID)

		randString := helper.RandomString(rand.Intn(80))
		noteTitle := "Title " + randString
		noteBody := "B"
		randCatId := rand.Intn(len(categoryList))
		noteCategory := categoryList[randCatId]

		requestBody := fmt.Sprintf(
			`{"title": "%s", "body": "%s", "id_category": %d}`,
			noteTitle, noteBody, noteCategory.ID)
		request := newTestRequest(updateUrl, http.MethodPut, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate web.ErrorResponse
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusBadRequest, responseCreate.Code)
		require.Equal(t, "BAD REQUEST", responseCreate.Status)
		require.Equal(t, "Body is should more than 2 characters", responseCreate.Message)
	})

	t.Run("Note_Update_BadRequest3_Fail", func(t *testing.T) {
		note := noteList[rand.Intn(len(noteList))]
		updateUrl := noteUrl + "/" + strconv.Itoa(note.ID)

		randString := helper.RandomString(rand.Intn(80))
		noteTitle := "Title " + randString
		noteBody := "Body " + randString
		noteCategory := 9999999

		requestBody := fmt.Sprintf(
			`{"title": "%s", "body": "%s", "id_category": %d}`,
			noteTitle, noteBody, noteCategory)
		request := newTestRequest(updateUrl, http.MethodPut, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate web.ErrorResponse
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusBadRequest, responseCreate.Code)
		require.Equal(t, "BAD REQUEST", responseCreate.Status)
		require.Equal(t, "category does not exists", responseCreate.Message)
	})
}

func TestDeleteNote(t *testing.T) {
	defer database.DeleteAllRecords(db)
	categoryList = database.CategorySeeder(db, 3)
	noteList := database.NoteSeeder(db, categoryList, 5)

	t.Run("Note_Delete_Success", func(t *testing.T) {
		note := noteList[rand.Intn(len(noteList))]
		deleteUrl := noteUrl + "/" + strconv.Itoa(note.ID)

		request := newTestRequest(deleteUrl, http.MethodDelete, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate testNoteJSON
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusOK, responseCreate.Code)
		require.Equal(t, "OK", responseCreate.Status)
		require.Equal(t, 0, responseCreate.Data.ID)
	})

	t.Run("Note_Delete_NotFound_Fail", func(t *testing.T) {
		deleteUrl := noteUrl + "/" + strconv.Itoa(9999999)

		request := newTestRequest(deleteUrl, http.MethodDelete, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate web.ErrorResponse
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusNotFound, responseCreate.Code)
		require.Equal(t, "NOT FOUND", responseCreate.Status)
		require.Equal(t, "note not found", responseCreate.Message)
	})
}
