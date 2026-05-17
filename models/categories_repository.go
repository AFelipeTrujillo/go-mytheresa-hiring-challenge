package models

import "gorm.io/gorm"

type CategoriesRepositoryInterface interface {
	GetAll() ([]Category, int, error)
	Create(category *Category) error
}

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) CategoriesRepositoryInterface {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) GetAll() ([]Category, int, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, len(categories), nil
}

func (r *CategoriesRepository) Create(category *Category) error {
	return r.db.Create(category).Error
}
