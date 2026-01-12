package user

import (
	"fmt"

	"github.com/NicoJCastro/gocourse_domain/domain"

	"log"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *domain.User) error
	GetAll(filters Filters, offset, limit int) ([]domain.User, error)
	Get(id string) (*domain.User, error)
	Delete(id string) error
	Update(id string, firstName *string, lastName *string, email *string, phone *string) error
	Count(filters Filters) (int64, error)
}

type repository struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepository(log *log.Logger, db *gorm.DB) Repository {
	return &repository{log: log, db: db}
}

func (r *repository) Create(user *domain.User) error {
	r.log.Println("---- Creating user in DB ----")
	result := r.db.Create(user)
	if result.Error != nil {
		r.log.Println("Error creating user: ", result.Error)
		return result.Error
	}
	r.log.Println("User created with ID: ", user.ID)
	return nil
}

func (r *repository) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {
	var users []domain.User
	tx := r.db.Model(&users)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&users)
	if result.Error != nil {
		r.log.Println("Error getting users: ", result.Error)
		return nil, result.Error
	}

	return users, nil

}

func (r *repository) Get(id string) (*domain.User, error) {
	user := domain.User{ID: id}
	result := r.db.First(&user)
	if result.Error != nil {
		r.log.Println("Error getting user: ", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func (r *repository) Delete(id string) error {
	user := domain.User{ID: id}
	result := r.db.Delete(&user)
	if result.Error != nil {
		r.log.Println("Error deleting user: ", result.Error)
		return result.Error
	}
	return nil
}

func (r *repository) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {

	updates := make(map[string]interface{})
	if firstName != nil {
		updates["first_name"] = *firstName
	}
	if lastName != nil {
		updates["last_name"] = *lastName
	}
	if email != nil {
		updates["email"] = *email
	}
	if phone != nil {
		updates["phone"] = *phone
	}
	result := r.db.Model(&domain.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		r.log.Println("Error updating user: ", result.Error)
		return result.Error
	}
	return nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("LOWER(first_name) LIKE ?", filters.FirstName)
	}

	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("LOWER(last_name) LIKE ?", filters.LastName)
	}

	if filters.Email != "" {
		filters.Email = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Email))
		tx = tx.Where("LOWER(email) LIKE ?", filters.Email)
	}

	if filters.Phone != "" {
		filters.Phone = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Phone))
		tx = tx.Where("LOWER(phone) LIKE ?", filters.Phone)
	}

	return tx
}

func (r *repository) Count(filters Filters) (int64, error) {
	var count int64
	tx := r.db.Model(&domain.User{})
	tx = applyFilters(tx, filters)
	result := tx.Count(&count)
	if result.Error != nil {
		r.log.Println("Error counting users: ", result.Error)
		return 0, result.Error
	}
	return count, nil
}
