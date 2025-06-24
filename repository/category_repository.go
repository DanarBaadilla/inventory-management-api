package repository

import (
	"inventory-management-api/model/domain"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll() ([]domain.Category, error)
	FindById(id int) (domain.Category, error)
	Save(category domain.Category) (domain.Category, error)
	Update(category domain.Category) (domain.Category, error)
	Delete(id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Order("id asc").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindById(id int) (domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	return category, err
}

func (r *categoryRepository) Save(category domain.Category) (domain.Category, error) {
	err := r.db.Create(&category).Error
	return category, err
}

func (r *categoryRepository) Update(category domain.Category) (domain.Category, error) {
	err := r.db.Model(&domain.Category{}).
		Where("id = ?", category.ID).
		Updates(map[string]interface{}{
			"name": category.Name,
		}).Error

	if err != nil {
		return domain.Category{}, err
	}

	var updatedCategory domain.Category
	err = r.db.First(&updatedCategory, category.ID).Error
	return updatedCategory, err
}

func (r *categoryRepository) Delete(id int) error {
	result := r.db.Delete(&domain.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
