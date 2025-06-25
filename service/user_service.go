package service

import (
	"errors"
	"fmt"
	"inventory-management-api/model/domain"
	"inventory-management-api/model/web"
	"inventory-management-api/repository"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindAll() ([]web.UserResponse, error)
	FindByID(id int) (web.UserResponse, error)
	Create(req web.UserCreateOrUpdateRequest) (web.UserResponse, error)
	Update(id int, req web.UserCreateOrUpdateRequest) (web.UserResponse, error)
	Delete(id int) error
}

type userServiceImpl struct {
	UserRepo repository.UserRepository
	Validate *validator.Validate
}

func NewUserService(userRepo repository.UserRepository, validate *validator.Validate) UserService {
	return &userServiceImpl{
		UserRepo: userRepo,
		Validate: validate,
	}
}

func (s *userServiceImpl) FindAll() ([]web.UserResponse, error) {
	users, err := s.UserRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []web.UserResponse
	for _, user := range users {
		responses = append(responses, toUserResponse(&user))
	}
	return responses, nil
}

func (s *userServiceImpl) FindByID(id int) (web.UserResponse, error) {
	user, err := s.UserRepo.FindByID(id)
	if err != nil {
		return web.UserResponse{}, err
	}
	return toUserResponse(user), nil
}

func (s *userServiceImpl) Create(req web.UserCreateOrUpdateRequest) (web.UserResponse, error) {
	// ✅ Validasi input
	if err := s.Validate.Struct(req); err != nil {
		return web.UserResponse{}, fmt.Errorf("validation error: %w", err)
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return web.UserResponse{}, err
	}

	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		Role:     req.Role,
	}

	createdUser, err := s.UserRepo.Save(&user)
	if err != nil {
		return web.UserResponse{}, err
	}

	return toUserResponse(createdUser), nil
}

func (s *userServiceImpl) Update(id int, req web.UserCreateOrUpdateRequest) (web.UserResponse, error) {
	// ✅ Validasi input
	if err := s.Validate.Struct(req); err != nil {
		return web.UserResponse{}, fmt.Errorf("validation error: %w", err)
	}

	user, err := s.UserRepo.FindByID(id)
	if err != nil {
		return web.UserResponse{}, errors.New("user not found")
	}

	// Hash password baru
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return web.UserResponse{}, err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Password = string(hashed)
	user.Role = req.Role

	updatedUser, err := s.UserRepo.Update(user)
	if err != nil {
		return web.UserResponse{}, err
	}

	return toUserResponse(updatedUser), nil
}

func (s *userServiceImpl) Delete(id int) error {
	user, err := s.UserRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	return s.UserRepo.Delete(user)
}

func toUserResponse(user *domain.User) web.UserResponse {
	return web.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}
