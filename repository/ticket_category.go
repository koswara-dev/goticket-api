package repository

import (
	"gotiket-api/model"

	"gorm.io/gorm"
)

type TicketCategoryRepository interface {
	Create(category *model.TicketCategory) error
	FindAll() ([]model.TicketCategory, error)
	FindByID(id int) (model.TicketCategory, error)
	Update(category *model.TicketCategory) error
	Delete(id int) error
}

type ticketCategoryRepository struct {
	db *gorm.DB
}

func NewTicketCategoryRepository(db *gorm.DB) TicketCategoryRepository {
	return &ticketCategoryRepository{db: db}
}

func (r *ticketCategoryRepository) Create(category *model.TicketCategory) error {
	return r.db.Create(category).Error
}

func (r *ticketCategoryRepository) FindAll() ([]model.TicketCategory, error) {
	var categories []model.TicketCategory
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *ticketCategoryRepository) FindByID(id int) (model.TicketCategory, error) {
	var category model.TicketCategory
	err := r.db.First(&category, id).Error
	return category, err
}

func (r *ticketCategoryRepository) Update(category *model.TicketCategory) error {
	return r.db.Save(category).Error
}

func (r *ticketCategoryRepository) Delete(id int) error {
	var category model.TicketCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return err
	}
	return r.db.Delete(&category).Error
}
