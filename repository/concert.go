package repository

import (
	"gotiket-api/model"

	"gorm.io/gorm"
)

type ConcertRepository interface {
	Create(concert *model.Concert) error
	FindByID(id uint) (model.Concert, error)
	FindByUserID(userID uint) (model.Concert, error)
	FindAll() ([]model.Concert, error)
	FindAllPagination(page int, limit int, search string, sort string) ([]model.Concert, int64, error)
	Update(concert *model.Concert) error
	Delete(id uint) error
}

type ConcertRepositoryImpl struct {
	db *gorm.DB
}

func NewConcertRepository(db *gorm.DB) ConcertRepository {
	return &ConcertRepositoryImpl{db: db}
}

func (r *ConcertRepositoryImpl) Create(concert *model.Concert) error {
	return r.db.Create(concert).Error
}

func (r *ConcertRepositoryImpl) FindByID(id uint) (model.Concert, error) {
	var concert model.Concert
	err := r.db.First(&concert, id).Error
	return concert, err
}

func (r *ConcertRepositoryImpl) FindByUserID(userID uint) (model.Concert, error) {
	var concert model.Concert
	err := r.db.Where("user_id = ?", userID).First(&concert).Error
	return concert, err
}

func (r *ConcertRepositoryImpl) FindAll() ([]model.Concert, error) {
	var concerts []model.Concert
	err := r.db.Find(&concerts).Error
	return concerts, err
}

// find all concert with pagination
func (r *ConcertRepositoryImpl) FindAllPagination(page int, limit int, search string, sort string) ([]model.Concert, int64, error) {
	var concerts []model.Concert
	var total int64

	query := r.db.Model(&model.Concert{})

	if search != "" {
		query = query.Where("title ILIKE ? OR venue ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch sort {
	case "title_asc":
		query = query.Order("title ASC")
	case "title_desc":
		query = query.Order("title DESC")
	case "date_asc":
		query = query.Order("date ASC")
	case "date_desc":
		query = query.Order("date DESC")
	default:
		query = query.Order("created_at DESC")
	}

	err := query.Limit(limit).Offset((page - 1) * limit).Find(&concerts).Error

	return concerts, total, err
}

func (r *ConcertRepositoryImpl) Update(concert *model.Concert) error {
	err := r.db.Save(concert).Error
	return err
}

func (r *ConcertRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.Concert{}, id).Error
	return err
}
