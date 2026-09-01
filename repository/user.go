package repository

import (
	"gorm.io/gorm"

	"gotiket-api/model"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByName(name string) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindByID(id uint) (model.User, error)
	Update(user *model.User) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(user *model.User) error {
	return r.db.Select("Name", "Email", "Password", "Role", "IsVerified").Create(user).Error
}

func (r *UserRepositoryImpl) FindByName(name string) (model.User, error) {
	var user model.User
	if err := r.db.Where("name = ?", name).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) FindByID(id uint) (model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) Update(user *model.User) error {
	return r.db.Save(user).Error
}
