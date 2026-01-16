package user

import (
	"context"
	"fmt"
	"log"

	"github.com/NicoJCastro/gocourse_domain/domain"
)

type (
	Filters struct {
		FirstName string
		LastName  string
		Email     string
		Phone     string
	}

	Service interface {
		Create(ctx context.Context, firstName, lastName, email, phone string) (*domain.User, error)
		Get(ctx context.Context, id string) (*domain.User, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) (*domain.User, error)
		Count(ctx context.Context, filters Filters) (int64, error)
	}
	// minúscula porque es privado
	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, firstName, lastName, email, phone string) (*domain.User, error) {
	s.log.Println("---- Creating user ----")

	// Validaciones básicas
	if firstName == "" || lastName == "" || email == "" || phone == "" {
		s.log.Println("Error: campos vacíos")
		return nil, fmt.Errorf("todos los campos son requeridos")
	}

	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	// Agregamos logging para debug
	s.log.Printf("Datos a insertar: %+v\n", user)

	// Propagamos el error del repositorio
	if err := s.repo.Create(ctx, &user); err != nil {
		s.log.Printf("Error creando usuario: %v\n", err)
		return nil, err
	}

	s.log.Printf("Usuario creado exitosamente: %s %s\n", firstName, lastName)
	return &user, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	s.log.Println("---- Getting all users ----")
	users, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		s.log.Printf("Error getting users: %v\n", err)
		return nil, err
	}
	return users, nil
}

func (s service) Get(ctx context.Context, id string) (*domain.User, error) {
	users, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Printf("Error getting user: %v\n", err)
		return nil, err
	}
	return users, nil
}

func (s service) Delete(ctx context.Context, id string) error {
	s.log.Println("---- Deleting user ----")
	return s.repo.Delete(ctx, id)
}

func (s service) Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) (*domain.User, error) {
	s.log.Println("---- Updating user ----")
	// ✅ Retornamos el usuario actualizado del repositorio
	user, err := s.repo.Update(ctx, id, firstName, lastName, email, phone)
	if err != nil {
		s.log.Printf("Error updating user: %v\n", err)
		return nil, err
	}
	s.log.Printf("User updated successfully: %s\n", id)
	return user, nil
}

func (s service) Count(ctx context.Context, filters Filters) (int64, error) {
	return s.repo.Count(ctx, filters)
}
