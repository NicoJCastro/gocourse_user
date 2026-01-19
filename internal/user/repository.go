package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/NicoJCastro/gocourse_domain/domain"

	"log"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) (*domain.User, error)
	Count(ctx context.Context, filters Filters) (int64, error)
}

type repository struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepository(log *log.Logger, db *gorm.DB) Repository {
	return &repository{log: log, db: db}
}

func (r *repository) Create(ctx context.Context, user *domain.User) error {
	r.log.Println("---- Creating user in DB ----")
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		r.log.Println("Error creating user: ", result.Error)
		return ErrUserNotCreated
	}
	r.log.Println("User created with ID: ", user.ID)
	return nil
}

func (r *repository) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	var users []domain.User
	tx := r.db.WithContext(ctx).Model(&users)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&users)
	if result.Error != nil {
		r.log.Println("Error getting users: ", result.Error)
		return nil, ErrUserNotRetrieved
	}

	return users, nil

}

func (r *repository) Get(ctx context.Context, id string) (*domain.User, error) {
	user := domain.User{ID: id}
	result := r.db.WithContext(ctx).First(&user)
	if result.Error != nil {
		r.log.Println("Error getting user: ", result.Error)
		// 🔍 Verificamos si es un error de GORM "record not found"
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, NewErrNotFound(id)
		}
		return nil, ErrUserNotRetrieved
	}
	return &user, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	user := domain.User{ID: id}
	result := r.db.WithContext(ctx).Delete(&user)
	if result.Error != nil {
		r.log.Println("Error deleting user: ", result.Error)
		// 🔍 Verificamos si es un error de GORM "record not found"
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return NewErrNotFound(id)
		}
		return ErrUserNotDeleted
	}
	if result.RowsAffected == 0 {
		r.log.Printf("No user found with ID: %s", id)
		return NewErrNotFound(id)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) (*domain.User, error) {
	// Construimos el mapa de updates solo con los campos proporcionados
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

	// Ejecutamos la actualización en la base de datos
	result := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		r.log.Println("Error updating user: ", result.Error)
		return nil, ErrUserNotUpdated
	}

	if result.RowsAffected == 0 {
		r.log.Printf("No user found with ID: %s", id)
		return nil, NewErrNotFound(id)
	}

	// Obtenemos el usuario actualizado después de la operación
	// Esto asegura que retornamos los datos más recientes (incluyendo timestamps)
	user, err := r.Get(ctx, id)
	if err != nil {
		r.log.Println("Error getting updated user: ", err)
		return nil, ErrUserNotRetrieved
	}

	return user, nil
}

func (r *repository) Count(ctx context.Context, filters Filters) (int64, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(&domain.User{})
	tx = applyFilters(tx, filters)
	result := tx.Count(&count)
	if result.Error != nil {
		r.log.Println("Error counting users: ", result.Error)
		return 0, ErrUserNotCounted
	}
	return count, nil
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
