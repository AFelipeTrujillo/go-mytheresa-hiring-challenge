package categories

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

type mockRepoCategories struct {
	categories []models.Category
	category   *models.Category
	total      int
	err        error
}

func (m *mockRepoCategories) GetAll() ([]models.Category, int, error) {
	return m.categories, m.total, m.err
}

func (m *mockRepoCategories) Create(category *models.Category) error {
	return m.err
}

func TestHandleGetAllCategories_Success(t *testing.T) {
	t.Run("json response with all categories", func(t *testing.T) {
		mock := &mockRepoCategories{
			categories: []models.Category{
				{Code: "CLOTHING", Name: "Clothing"},
				{Code: "SHOES", Name: "Shoes"},
				{Code: "ACCESSORIES", Name: "Accessories"},
			},
			total: 3,
		}

		handler := NewCategoriesHandler(mock)

		req := httptest.NewRequest("GET", "/categories", nil)
		rec := httptest.NewRecorder()

		handler.HandleGet(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Status 200 was expected, but %d was returned", rec.Code)
		}

		assert.Equal(t, http.StatusOK, rec.Code)

		var response CategoryResponseDTO

		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, "CLOTHING", response.Categories[0].Code)
		assert.Equal(t, "SHOES", response.Categories[1].Code)
		assert.Equal(t, "ACCESSORIES", response.Categories[2].Code)

		assert.Equal(t, 3, response.Total)

	})
}

func TestHandleCreateCategories_Success(t *testing.T) {
	t.Run("JSON response with a created category", func(t *testing.T) {
		mock := &mockRepoCategories{}

		handler := NewCategoriesHandler(mock)
		body := `{"code":"NEWCAT","name":"New Category"}`
		req := httptest.NewRequest("POST", "/categories", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.HandleCreateCategory(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Status 200 was expected, but %d was returned", rec.Code)
		}

		assert.Equal(t, http.StatusOK, rec.Code)
		var category models.Category
		json.Unmarshal(rec.Body.Bytes(), &category)

		assert.Equal(t, "NEWCAT", category.Code)
		assert.Equal(t, "New Category", category.Name)
	})
}

func TestHandleCreateCategory_InvalidJSON(t *testing.T) {
	t.Run("returns 400 for malformed JSON", func(t *testing.T) {
		mock := &mockRepoCategories{}
		handler := NewCategoriesHandler(mock)

		body := `{invalid}`
		req := httptest.NewRequest("POST", "/categories", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.HandleCreateCategory(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandleCreateCategory_ValidationError(t *testing.T) {
	t.Run("returns 400 for missing required fields", func(t *testing.T) {
		mock := &mockRepoCategories{}
		handler := NewCategoriesHandler(mock)

		body := `{"code":"","name":"Test"}`
		req := httptest.NewRequest("POST", "/categories", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.HandleCreateCategory(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
