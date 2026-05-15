package catalog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockRepo struct {
	products []models.Product
	product  *models.Product
	total    int
	err      error
}

func (m *mockRepo) GetAllProducts(offset, limit int, category string, priceLessThan float64) ([]models.Product, int, error) {
	return m.products, m.total, m.err
}

func (m *mockRepo) GetProductByCode(code string) (*models.Product, error) {
	return m.product, m.err
}

func TestHandleGetSuccess(t *testing.T) {
	mock := &mockRepo{
		products: []models.Product{
			{
				Code:  "PROD001",
				Price: decimal.NewFromFloat(9.99),
				Category: models.Category{
					Code: "CLOTHING",
					Name: "Clothing",
				},
			},
		},
		total: 1,
	}

	handler := NewCatalogHandler(mock)
	req := httptest.NewRequest("GET", "/catalog", nil)
	rec := httptest.NewRecorder()

	handler.HandleGet(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response Response
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, 1, response.Total)
	assert.Equal(t, "PROD001", response.Products[0].Code)
	assert.Equal(t, "Clothing", response.Products[0].Category)

}

func TestHandleGetByCode_Success(t *testing.T) {

	t.Run("json response with product by code", func(t *testing.T) {
		mock := &mockRepo{
			product: &models.Product{
				Code:  "PROD001",
				Price: decimal.NewFromFloat(6.66),
				Category: models.Category{
					Code: "CLOTHING",
					Name: "Clothing",
				},
				Variants: []models.Variant{
					{Name: "Variant A", SKU: "SKU001A", Price: decimal.NewFromFloat(11.99)},
					{Name: "Variant B", SKU: "SKU001B", Price: decimal.Decimal{}},
				},
			},
		}

		handler := NewCatalogHandler(mock)
		req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
		rec := httptest.NewRecorder()

		handler.HandleGetByCode(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Status 200 was expected, but %d was returned", rec.Code)
		}

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ProductDetailResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, "PROD001", response.Code)
		assert.Equal(t, 6.66, response.Price)
		assert.Equal(t, "Clothing", response.Category)
		assert.Len(t, response.Variants, 2)

		assert.Equal(t, 11.99, response.Variants[0].Price)
		assert.Equal(t, 6.66, response.Variants[1].Price)

	})

}

func TestHandleGetByCode_NotFound(t *testing.T) {
	t.Run("json respnse with product not found", func(t *testing.T) {
		mock := &mockRepo{
			err: gorm.ErrRecordNotFound,
		}

		handlers := NewCatalogHandler(mock)
		req := httptest.NewRequest("GET", "/catalog/NOTFOUND", nil)
		rec := httptest.NewRecorder()

		handlers.HandleGetByCode(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
