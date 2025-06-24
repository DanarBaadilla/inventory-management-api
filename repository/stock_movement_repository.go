package repository

import (
	"inventory-management-api/model/domain"

	"gorm.io/gorm"
)

type StockMovementRepository interface {
	FindAll() ([]domain.StockMovement, error)
	FindById(id int) (domain.StockMovement, error)
	Save(movement domain.StockMovement, tx *gorm.DB) (domain.StockMovement, error)
	Delete(id int) error
	FindByMonth(month string, filters map[string]interface{}) ([]domain.StockMovement, error)
}

type stockMovementRepository struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) StockMovementRepository {
	return &stockMovementRepository{db: db}
}

func (r *stockMovementRepository) FindAll() ([]domain.StockMovement, error) {
	var movements []domain.StockMovement
	err := r.db.Order("id desc").Find(&movements).Error
	return movements, err
}

func (r *stockMovementRepository) FindById(id int) (domain.StockMovement, error) {
	var m domain.StockMovement
	err := r.db.First(&m, id).Error
	return m, err
}

func (r *stockMovementRepository) Save(m domain.StockMovement, tx *gorm.DB) (domain.StockMovement, error) {
	err := tx.Create(&m).Error
	return m, err
}

func (r *stockMovementRepository) Delete(id int) error {
	return r.db.Delete(&domain.StockMovement{}, id).Error
}

// ✅ Fleksibel: Jika month kosong, maka tidak difilter berdasarkan bulan
func (r *stockMovementRepository) FindByMonth(month string, filters map[string]interface{}) ([]domain.StockMovement, error) {
	query := r.db

	if month != "" {
		query = query.Where("DATE_FORMAT(created_at, '%Y-%m') = ?", month)
	}

	if productID, ok := filters["product_id"]; ok {
		query = query.Where("product_id = ?", productID)
	}
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if movementType, ok := filters["type"]; ok {
		query = query.Where("type = ?", movementType)
	}

	var movements []domain.StockMovement
	err := query.Order("created_at desc").Find(&movements).Error
	return movements, err
}
