package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
}

type Product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type CatalogHandler struct {
	repo models.ProductsRepositoryInterface
}

func NewCatalogHandler(r models.ProductsRepositoryInterface) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	if limit < 1 {
		limit = 1
	}

	if limit > 100 {
		limit = 100
	}

	category := r.URL.Query().Get("category")

	priceLessThanStr := r.URL.Query().Get("price_less_than")
	var priceLessThan float64
	if priceLessThanStr != "" {
		priceLessThan, err = strconv.ParseFloat(priceLessThanStr, 64)
		if err != nil || priceLessThan < 0 {
			http.Error(w, "invalid price_less_than", http.StatusBadRequest)
			return
		}
	}

	res, total, err := h.repo.GetAllProducts(offset, limit, category, priceLessThan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products: products,
		Total:    total,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
