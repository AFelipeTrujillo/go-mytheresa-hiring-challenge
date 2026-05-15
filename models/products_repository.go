package models

import (
	"gorm.io/gorm"
)

type ProductsRepositoryInterface interface {
	GetAllProducts(offset, limit int) ([]Product, int, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductsRepositoryInterface {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(offset, limit int) ([]Product, int, error) {
	var products []Product
	var total int64

	if err := r.db.Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Variants").Preload("Category").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, int(total), nil
}
