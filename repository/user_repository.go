package repository

import (
	"inventory-management-api/model/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	FindByID(id int) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Save(user *domain.User) (*domain.User, error)
	Update(user *domain.User) (*domain.User, error)
	Delete(user *domain.User) error
}

type userRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (r *userRepositoryImpl) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepositoryImpl) FindByID(id int) (*domain.User, error) {
	var user domain.User
	err := r.DB.First(&user, id).Error
	return &user, err
}

func (r *userRepositoryImpl) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *userRepositoryImpl) Save(user *domain.User) (*domain.User, error) {
	err := r.DB.Create(user).Error
	return user, err
}

func (r *userRepositoryImpl) Update(user *domain.User) (*domain.User, error) {
	err := r.DB.Save(user).Error
	return user, err
}

func (r *userRepositoryImpl) Delete(user *domain.User) error {
	return r.DB.Delete(user).Error
}
