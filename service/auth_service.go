package service

import (
	"errors"
	"inventory-management-api/helper"
	"inventory-management-api/model/web"
	"inventory-management-api/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email string, password string) (web.LoginResponse, error)
}

type authServiceImpl struct {
	UserRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authServiceImpl{UserRepo: userRepo}
}

func (s *authServiceImpl) Login(email, password string) (web.LoginResponse, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return web.LoginResponse{}, errors.New("email or password is incorrect")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return web.LoginResponse{}, errors.New("email or password is incorrect")
	}

	token, err := helper.GenerateToken(user.ID, user.Role)
	if err != nil {
		return web.LoginResponse{}, errors.New("failed to generate token")
	}

	return web.LoginResponse{
		Token: token,
		User: web.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
