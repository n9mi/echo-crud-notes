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
	"github.com/naomigrain/echo-crud-notes/model/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCategoryJSON struct {
	Code   int              `json:"code"`
	Status string           `json:"status"`
	Data   web.CategoryJSON `json:"data"`
}

type testCategoryListJSON struct {
	Code   int                `json:"code"`
	Status string             `json:"status"`
	Data   []web.CategoryJSON `json:"data"`
}

var categoryUrl string = "http://127.0.0.1:8000/api/categories"

func TestCreateCategory(t *testing.T) {
	defer database.DeleteCategoryRecords(db)
	createUrl := categoryUrl
	t.Run("Category_Create_Success", func(t *testing.T) {
		categoryName := "Category " + helper.RandomString(10)
		requestBody := fmt.Sprintf(
			`{"name": "%s"}`, categoryName)
		request := newTestRequest(createUrl,
			http.MethodPost, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate testCategoryJSON
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusOK, responseCreate.Code)
		require.Equal(t, "OK", responseCreate.Status)
		require.Equal(t, categoryName, responseCreate.Data.Name)
	})
	t.Run("Category_Create_Validation_Fail", func(t *testing.T) {
		categoryName := helper.RandomString(1)
		requestBody := fmt.Sprintf(
			`{"name": "%s"}`, categoryName)
		request := newTestRequest(createUrl,
			http.MethodPost, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseCreate web.ErrorResponse
		json.Unmarshal(responseBody, &responseCreate)

		require.Equal(t, http.StatusBadRequest, responseCreate.Code)
		require.Equal(t, "BAD REQUEST", responseCreate.Status)
		require.Equal(t, "Name is should more than 2 characters", responseCreate.Message)
	})
}

func TestUpdateCategory(t *testing.T) {
	defer database.DeleteCategoryRecords(db)
	categories := database.CategorySeeder(db, 2)
	t.Run("Category_Update_Success", func(t *testing.T) {
		updateUrl := categoryUrl + "/" + strconv.Itoa(categories[0].ID)
		categoryNameUpdate := "Category " + helper.RandomString(rand.Intn(80))
		requestBody := fmt.Sprintf(
			`{"name": "%s"}`, categoryNameUpdate,
		)
		request := newTestRequest(updateUrl, http.MethodPut, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseUpdate testCategoryJSON
		json.Unmarshal(responseBody, &responseUpdate)

		require.Equal(t, http.StatusOK, responseUpdate.Code)
		require.Equal(t, "OK", responseUpdate.Status)
		require.Equal(t, categories[0].ID, responseUpdate.Data.ID)
		require.Equal(t, categoryNameUpdate, responseUpdate.Data.Name)
	})
	t.Run("Category_Update_Validation_Fail", func(t *testing.T) {
		updateUrl := categoryUrl + "/" + strconv.Itoa(categories[1].ID)
		categoryNameUpdate := helper.RandomString(1)
		requestBody := fmt.Sprintf(
			`{"name": "%s"}`, categoryNameUpdate,
		)
		request := newTestRequest(updateUrl, http.MethodPut, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseUpdate web.ErrorResponse
		json.Unmarshal(responseBody, &responseUpdate)

		require.Equal(t, http.StatusBadRequest, responseUpdate.Code)
		require.Equal(t, "BAD REQUEST", responseUpdate.Status)
		require.Equal(t, "Name is should more than 2 characters", responseUpdate.Message)
	})
	t.Run("Category_Update_NotFound_Fail", func(t *testing.T) {
		updateUrl := categoryUrl + "/" + strconv.Itoa(9999999)
		categoryNameUpdate := "Category " + helper.RandomString(rand.Intn(80))
		requestBody := fmt.Sprintf(
			`{"name": "%s"}`, categoryNameUpdate,
		)
		request := newTestRequest(updateUrl, http.MethodPut, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseUpdate web.ErrorResponse
		json.Unmarshal(responseBody, &responseUpdate)

		require.Equal(t, http.StatusNotFound, responseUpdate.Code)
		require.Equal(t, "NOT FOUND", responseUpdate.Status)
		require.Equal(t, "category not found", responseUpdate.Message)
	})
	t.Run("Category_Update_NotFound2_Fail", func(t *testing.T) {
		updateUrl := categoryUrl + "/thisIsID"
		categoryNameUpdate := "Category " + helper.RandomString(rand.Intn(80))
		requestBody := fmt.Sprintf(
			`{"name": "%s"}`, categoryNameUpdate,
		)
		request := newTestRequest(updateUrl, http.MethodPut, requestBody)

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseUpdate web.ErrorResponse
		json.Unmarshal(responseBody, &responseUpdate)

		require.Equal(t, http.StatusNotFound, responseUpdate.Code)
		require.Equal(t, "NOT FOUND", responseUpdate.Status)
		require.Equal(t, "category not found", responseUpdate.Message)
	})
}

func TestDeleteCategory(t *testing.T) {
	defer database.DeleteCategoryRecords(db)
	categories := database.CategorySeeder(db, 1)

	t.Run("Category_Delete_Success", func(t *testing.T) {
		deleteUrl := categoryUrl + "/" + strconv.Itoa(categories[0].ID)
		request := newTestRequest(deleteUrl, http.MethodDelete, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseDelete testCategoryJSON
		json.Unmarshal(responseBody, &responseDelete)

		require.Equal(t, http.StatusOK, responseDelete.Code)
		require.Equal(t, "OK", responseDelete.Status)
	})
	t.Run("Category_Delete_NotFound_Fail", func(t *testing.T) {
		deleteUrl := categoryUrl + "/" + strconv.Itoa(9999999)
		request := newTestRequest(deleteUrl, http.MethodDelete, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseDelete web.ErrorResponse
		json.Unmarshal(responseBody, &responseDelete)

		require.Equal(t, http.StatusNotFound, responseDelete.Code)
		require.Equal(t, "NOT FOUND", responseDelete.Status)
		require.Equal(t, "category not found", responseDelete.Message)
	})
	t.Run("Category_Delete_NotFound2_Fail", func(t *testing.T) {
		deleteUrl := categoryUrl + "/thisistestforid"
		request := newTestRequest(deleteUrl, http.MethodDelete, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseDelete web.ErrorResponse
		json.Unmarshal(responseBody, &responseDelete)

		require.Equal(t, http.StatusNotFound, responseDelete.Code)
		require.Equal(t, "NOT FOUND", responseDelete.Status)
		require.Equal(t, "category not found", responseDelete.Message)
	})
}

func TestGetCategories(t *testing.T) {
	defer database.DeleteCategoryRecords(db)
	categories := database.CategorySeeder(db, 3)

	t.Run("Category_FindById_Success", func(t *testing.T) {
		getByIdUrl := categoryUrl + "/" + strconv.Itoa(categories[0].ID)
		request := newTestRequest(getByIdUrl, http.MethodGet, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseGetById testCategoryJSON
		json.Unmarshal(responseBody, &responseGetById)

		require.Equal(t, http.StatusOK, responseGetById.Code)
		require.Equal(t, "OK", responseGetById.Status)
		require.Equal(t, categories[0].ID, responseGetById.Data.ID)
		require.Equal(t, categories[0].Name, responseGetById.Data.Name)
	})
	t.Run("Category_FindById_NotFound_Fail", func(t *testing.T) {
		getByIdUrl := categoryUrl + "/" + strconv.Itoa(9999999)
		request := newTestRequest(getByIdUrl, http.MethodGet, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseGetById web.ErrorResponse
		json.Unmarshal(responseBody, &responseGetById)

		require.Equal(t, http.StatusNotFound, responseGetById.Code)
		require.Equal(t, "NOT FOUND", responseGetById.Status)
		require.Equal(t, "category not found", responseGetById.Message)
	})
	t.Run("Category_FindById_NotFound2_Fail", func(t *testing.T) {
		getByIdUrl := categoryUrl + "/thisisidok"
		request := newTestRequest(getByIdUrl, http.MethodGet, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseGetById web.ErrorResponse
		json.Unmarshal(responseBody, &responseGetById)

		require.Equal(t, http.StatusNotFound, responseGetById.Code)
		require.Equal(t, "NOT FOUND", responseGetById.Status)
		require.Equal(t, "category not found", responseGetById.Message)
	})
	t.Run("Category_GetAll_Success", func(t *testing.T) {
		request := newTestRequest(categoryUrl, http.MethodGet, "")

		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		response := recorder.Result()

		responseBody, _ := io.ReadAll(response.Body)
		var responseGetAll testCategoryListJSON
		json.Unmarshal(responseBody, &responseGetAll)

		require.Equal(t, http.StatusOK, responseGetAll.Code)
		require.Equal(t, "OK", responseGetAll.Status)
		require.Equal(t, len(categories), len(responseGetAll.Data))
		for i := 0; i < len(categories); i++ {
			assert.Equal(t, categories[i].ID, responseGetAll.Data[i].ID)
			assert.Equal(t, categories[i].Name, responseGetAll.Data[i].Name)
		}
	})
}

func TestDeleteCategories(t *testing.T) {
	database.DeleteAllRecords(db)
}
