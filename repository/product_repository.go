package repository

import (
	"errors"
	"inventory-management-api/model/domain"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll() ([]domain.Product, error)
	FindById(id int) (domain.Product, error)
	Save(product domain.Product) (domain.Product, error)
	Update(product domain.Product) (domain.Product, error)
	Delete(id int) error
	SearchWithFilter(name, sort string, page, limit int) ([]domain.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Order("id asc").Find(&products).Error
	return products, err
}

func (r *productRepository) FindById(id int) (domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Product{}, errors.New("product not found")
	}
	return product, err
}

func (r *productRepository) Save(product domain.Product) (domain.Product, error) {
	err := r.db.Create(&product).Error
	return product, err
}

func (r *productRepository) Update(product domain.Product) (domain.Product, error) {
	// Ensure product exists before update
	var existing domain.Product
	err := r.db.First(&existing, product.ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Product{}, errors.New("product not found")
	}

	err = r.db.Model(&existing).Updates(product).Error
	return existing, err
}

func (r *productRepository) Delete(id int) error {
	result := r.db.Delete(&domain.Product{}, id)
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return result.Error
}

func (r *productRepository) SearchWithFilter(name, sort string, page, limit int) ([]domain.Product, error) {
	var products []domain.Product
	db := r.db.Model(&domain.Product{})

	// Filter
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	// Sorting
	switch sort {
	case "name_asc":
		db = db.Order("name ASC")
	case "name_desc":
		db = db.Order("name DESC")
	case "stock_asc":
		db = db.Order("stock ASC")
	case "stock_desc":
		db = db.Order("stock DESC")
	case "created_asc":
		db = db.Order("created_at ASC")
	case "created_desc":
		db = db.Order("created_at DESC")
	default:
		db = db.Order("id ASC")
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit > 0 {
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
	}

	err := db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, errors.New("no product found")
	}
	return products, nil
}
