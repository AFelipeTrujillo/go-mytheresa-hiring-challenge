package models

import (
	"gorm.io/gorm"
)

type ProductsRepositoryInterface interface {
	GetAllProducts(offset, limit int, category string, priceLessThan float64) ([]Product, int, error)
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
