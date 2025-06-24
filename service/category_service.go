package service

import (
	"fmt"
	"inventory-management-api/model/domain"
	"inventory-management-api/model/web"
	"inventory-management-api/repository"

	"github.com/go-playground/validator/v10"
)

type CategoryService interface {
	FindAll() ([]web.CategoryResponse, error)
	FindById(id int) (web.CategoryResponse, error)
	Create(request web.CategoryCreateOrUpdateRequest) (web.CategoryResponse, error)
	Update(id int, request web.CategoryCreateOrUpdateRequest) (web.CategoryResponse, error)
	Delete(id int) error
}

type categoryService struct {
	Repository repository.CategoryRepository
	Validate   *validator.Validate
}

func NewCategoryService(repo repository.CategoryRepository, validate *validator.Validate) CategoryService {
	return &categoryService{
		Repository: repo,
		Validate:   validate,
	}
}

func (s *categoryService) FindAll() ([]web.CategoryResponse, error) {
	categories, err := s.Repository.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []web.CategoryResponse
	for _, c := range categories {
		responses = append(responses, web.CategoryResponse{
			ID:   c.ID,
			Name: c.Name,
		})
	}
	return responses, nil
}

func (s *categoryService) FindById(id int) (web.CategoryResponse, error) {
	c, err := s.Repository.FindById(id)
	if err != nil {
		return web.CategoryResponse{}, err
	}
	return web.CategoryResponse{ID: c.ID, Name: c.Name}, nil
}

func (s *categoryService) Create(req web.CategoryCreateOrUpdateRequest) (web.CategoryResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return web.CategoryResponse{}, fmt.Errorf("validation error: %w", err)
	}

	category := domain.Category{
		Name: req.Name,
	}
	saved, err := s.Repository.Save(category)
	if err != nil {
		return web.CategoryResponse{}, err
	}
	return web.CategoryResponse{ID: saved.ID, Name: saved.Name}, nil
}

func (s *categoryService) Update(id int, req web.CategoryCreateOrUpdateRequest) (web.CategoryResponse, error) {
	// Validasi input
	if err := s.Validate.Struct(req); err != nil {
		return web.CategoryResponse{}, fmt.Errorf("validation error: %w", err)
	}

	// Ambil data lama dari database
	existingCategory, err := s.Repository.FindById(id)
	if err != nil {
		return web.CategoryResponse{}, err
	}

	// Ubah hanya field yang diizinkan
	existingCategory.Name = req.Name

	// Simpan kembali ke database
	updated, err := s.Repository.Update(existingCategory)
	if err != nil {
		return web.CategoryResponse{}, err
	}

	// Return response
	return web.CategoryResponse{ID: updated.ID, Name: updated.Name}, nil
}

func (s *categoryService) Delete(id int) error {
	return s.Repository.Delete(id)
}
