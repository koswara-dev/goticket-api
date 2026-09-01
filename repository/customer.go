package repository

import (
	"gotiket-api/model"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *model.Customer) error
	FindByID(id uint) (model.Customer, error)
	FindByUserID(userID uint) (model.Customer, error)
	FindAll() ([]model.Customer, error)
	FindAllPagination(page int, limit int, search string, sort string) ([]model.Customer, int64, error)
	Update(customer *model.Customer) error
	Delete(id uint) error
}

type CustomerRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &CustomerRepositoryImpl{db: db}
}

func (r *CustomerRepositoryImpl) Create(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepositoryImpl) FindByID(id uint) (model.Customer, error) {
	var customer model.Customer
	err := r.db.First(&customer, id).Error
	return customer, err
}

func (r *CustomerRepositoryImpl) FindByUserID(userID uint) (model.Customer, error) {
	var customer model.Customer
	err := r.db.Where("user_id = ?", userID).First(&customer).Error
	return customer, err
}

func (r *CustomerRepositoryImpl) FindAll() ([]model.Customer, error) {
	var customers []model.Customer
	err := r.db.Find(&customers).Error
	return customers, err
}

// find all customer with pagination
func (r *CustomerRepositoryImpl) FindAllPagination(page int, limit int, search string, sort string) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var total int64

	query := r.db.Model(&model.Customer{})

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ? OR nik ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch sort {
	case "name_asc":
		query = query.Order("name ASC")
	case "name_desc":
		query = query.Order("name DESC")
	case "email_asc":
		query = query.Order("email ASC")
	case "email_desc":
		query = query.Order("email DESC")
	case "nik_asc":
		query = query.Order("nik ASC")
	case "nik_desc":
		query = query.Order("nik DESC")
	default:
		query = query.Order("created_at DESC")
	}

	err := query.Limit(limit).Offset((page - 1) * limit).Find(&customers).Error

	return customers, total, err
}

func (r *CustomerRepositoryImpl) Update(customer *model.Customer) error {
	err := r.db.Save(customer).Error
	return err
}

func (r *CustomerRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.Customer{}, id).Error
	return err
}
