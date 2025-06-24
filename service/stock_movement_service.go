package service

import (
	"errors"
	"inventory-management-api/model/domain"
	"inventory-management-api/model/web"
	"inventory-management-api/repository"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type StockMovementService interface {
	FindAll() ([]web.StockMovementResponse, error)
	FindById(id int) (web.StockMovementResponse, error)
	Create(userID int, req web.StockMovementCreateRequest) (web.StockMovementResponse, error)
	Delete(id int) error
	GetMonthlyReport(month string, filters map[string]interface{}) ([]web.StockMovementResponse, error)
}

type stockMovementService struct {
	RepoMovement repository.StockMovementRepository
	RepoProduct  repository.ProductRepository
	DB           *gorm.DB
	Validate     *validator.Validate
}

func NewStockMovementService(
	repoMovement repository.StockMovementRepository,
	repoProduct repository.ProductRepository,
	db *gorm.DB,
	validate *validator.Validate,
) StockMovementService {
	return &stockMovementService{
		RepoMovement: repoMovement,
		RepoProduct:  repoProduct,
		DB:           db,
		Validate:     validate,
	}
}

func (s *stockMovementService) FindAll() ([]web.StockMovementResponse, error) {
	movements, err := s.RepoMovement.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []web.StockMovementResponse
	for _, m := range movements {
		responses = append(responses, toStockMovementResponse(m))
	}
	return responses, nil
}

func (s *stockMovementService) FindById(id int) (web.StockMovementResponse, error) {
	m, err := s.RepoMovement.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return web.StockMovementResponse{}, errors.New("stock movement not found")
		}
		return web.StockMovementResponse{}, err
	}
	return toStockMovementResponse(m), nil
}

func (s *stockMovementService) Create(userID int, req web.StockMovementCreateRequest) (web.StockMovementResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return web.StockMovementResponse{}, errors.New("validation failed")
	}

	tx := s.DB.Begin()

	// Cari produk
	product, err := s.RepoProduct.FindById(req.ProductID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return web.StockMovementResponse{}, errors.New("product not found")
		}
		return web.StockMovementResponse{}, err
	}

	// Validasi dan update stok
	if req.Type == "in" {
		product.Stock += req.Quantity
	} else if req.Type == "out" {
		if product.Stock < req.Quantity {
			tx.Rollback()
			return web.StockMovementResponse{}, errors.New("stock not enough")
		}
		product.Stock -= req.Quantity
	}

	// Simpan update stock product
	_, err = s.RepoProduct.Update(product)
	if err != nil {
		tx.Rollback()
		return web.StockMovementResponse{}, err
	}

	// Simpan movement
	movement := domain.StockMovement{
		ProductID: req.ProductID,
		UserID:    userID,
		Type:      req.Type,
		Quantity:  req.Quantity,
		Note:      req.Note,
	}

	saved, err := s.RepoMovement.Save(movement, tx)
	if err != nil {
		tx.Rollback()
		return web.StockMovementResponse{}, err
	}

	tx.Commit()
	return toStockMovementResponse(saved), nil
}

func (s *stockMovementService) Delete(id int) error {
	_, err := s.RepoMovement.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("stock movement not found")
		}
		return err
	}

	return s.RepoMovement.Delete(id)
}

func (s *stockMovementService) GetMonthlyReport(month string, filters map[string]interface{}) ([]web.StockMovementResponse, error) {
	movements, err := s.RepoMovement.FindByMonth(month, filters)
	if err != nil {
		return nil, err
	}

	// ✅ Kembalikan error jika tidak ada data
	if len(movements) == 0 {
		return nil, errors.New("report not found")
	}

	var responses []web.StockMovementResponse
	for _, m := range movements {
		responses = append(responses, toStockMovementResponse(m))
	}
	return responses, nil
}

func toStockMovementResponse(m domain.StockMovement) web.StockMovementResponse {
	return web.StockMovementResponse{
		ID:        m.ID,
		ProductID: m.ProductID,
		UserID:    m.UserID,
		Type:      m.Type,
		Quantity:  m.Quantity,
		Note:      m.Note,
		CreatedAt: m.CreatedAt,
	}
}
