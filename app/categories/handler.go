package categories

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoryResponseDTO struct {
	Categories []CategoryDTO `json:"categories"`
	Total      int           `json:"total"`
}

type CreateCategoryRequestDTO struct {
	Code string `json:"code" validate:"required,alpha,min=3,max=32"`
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type CategoryDTO struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoriesHandler struct {
	repo models.CategoriesRepositoryInterface
}

var validate = validator.New()

func NewCategoriesHandler(r models.CategoriesRepositoryInterface) *CategoriesHandler {
	return &CategoriesHandler{repo: r}
}

func (h *CategoriesHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	res, total, err := h.repo.GetAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	categories := make([]CategoryDTO, len(res))
	for i, c := range res {
		categories[i] = CategoryDTO{
			Code: c.Code,
			Name: c.Name,
		}
	}

	w.Header().Set("Content-Type", "application/json")

	response := CategoryResponseDTO{
		Categories: categories,
		Total:      total,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (h *CategoriesHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid JSON payload or bad data types")
		return
	}

	if err := validate.Struct(req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}

	category := models.Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := h.repo.Create(&category); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, category)

}
