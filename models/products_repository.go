package models

import (
	"gorm.io/gorm"
)

type ProductsRepositoryInterface interface {
	GetAllProducts(offset, limit int, category string, priceLessThan float64) ([]Product, int, error)
	GetProductByCode(code string) (*Product, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductsRepositoryInterface {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(offset, limit int, category string, priceLessThan float64) ([]Product, int, error) {
	var products []Product
	var total int64

	query := r.db.Model(&Product{})

	if category != "" {
		query = query.Joins("Category").Where("Category.code = ?", category)
	}

	if priceLessThan > 0 {
		query = query.Where("products.price < ?", priceLessThan)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Variants").Preload("Category").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, int(total), nil
}

func (r *ProductsRepository) GetProductByCode(code string) (*Product, error) {

	var product Product

	query := r.db.Model(&Product{})

	query = query.Where("code = ?", code)

	if err := query.Preload("Variants").Preload("Category").First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}
