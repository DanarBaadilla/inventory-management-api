package service

import (
	"errors"
	"fmt"
	"inventory-management-api/model/domain"
	"inventory-management-api/model/web"
	"inventory-management-api/repository"

	"github.com/go-playground/validator/v10"
)

type ProductService interface {
	FindAll() ([]web.ProductResponse, error)
	FindById(id int) (web.ProductResponse, error)
	Create(request web.ProductCreateOrUpdateRequest) (web.ProductResponse, error)
	Update(id int, request web.ProductCreateOrUpdateRequest) (web.ProductResponse, error)
	Delete(id int) error
	SearchWithFilter(name, sort string, page, limit int) ([]web.ProductResponse, error)
}

type productService struct {
	Repo     repository.ProductRepository
	Validate *validator.Validate
}

func NewProductService(repo repository.ProductRepository, validate *validator.Validate) ProductService {
	return &productService{
		Repo:     repo,
		Validate: validate,
	}
}

func (s *productService) FindAll() ([]web.ProductResponse, error) {
	products, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}
	return toProductResponses(products), nil
}

func (s *productService) FindById(id int) (web.ProductResponse, error) {
	p, err := s.Repo.FindById(id)
	if err != nil {
		return web.ProductResponse{}, fmt.Errorf("product not found")
	}
	return toProductResponse(p), nil
}

func (s *productService) Create(req web.ProductCreateOrUpdateRequest) (web.ProductResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return web.ProductResponse{}, fmt.Errorf("validation error: %w", err)
	}

	product := domain.Product{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Stock:      req.Stock,
	}

	saved, err := s.Repo.Save(product)
	if err != nil {
		return web.ProductResponse{}, err
	}
	return toProductResponse(saved), nil
}

func (s *productService) Update(id int, req web.ProductCreateOrUpdateRequest) (web.ProductResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return web.ProductResponse{}, fmt.Errorf("validation error: %w", err)
	}

	// pastikan data lama ada
	_, err := s.Repo.FindById(id)
	if err != nil {
		return web.ProductResponse{}, errors.New("product not found")
	}

	product := domain.Product{
		ID:         id,
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Stock:      req.Stock,
	}

	updated, err := s.Repo.Update(product)
	if err != nil {
		return web.ProductResponse{}, err
	}

	return toProductResponse(updated), nil
}

func (s *productService) Delete(id int) error {
	_, err := s.Repo.FindById(id)
	if err != nil {
		return errors.New("product not found")
	}
	return s.Repo.Delete(id)
}

func (s *productService) SearchWithFilter(name, sort string, page, limit int) ([]web.ProductResponse, error) {
	products, err := s.Repo.SearchWithFilter(name, sort, page, limit)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, errors.New("no products found")
	}
	return toProductResponses(products), nil
}

// Helpers

func toProductResponse(p domain.Product) web.ProductResponse {
	return web.ProductResponse{
		ID:         p.ID,
		Name:       p.Name,
		CategoryID: p.CategoryID,
		Stock:      p.Stock,
	}
}

func toProductResponses(products []domain.Product) []web.ProductResponse {
	var responses []web.ProductResponse
	for _, p := range products {
		responses = append(responses, toProductResponse(p))
	}
	return responses
}
