package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/NicoJCastro/gocourse_meta/meta"
)

type (
	Controller func(ctx context.Context, request interface{}) (response interface{}, err error)

	Endpoint struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	GetRequest struct {
		ID string `json:"id"`
	}

	DeleteRequest struct {
		ID string `json:"id"`
	}

	GetAllRequest struct {
		FirstName string
		LastName  string
		Email     string
		Phone     string
		Limit     int
		Page      int
	}

	UpdateRequest struct {
		ID        string  `json:"id"`
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoint {
	return Endpoint{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(CreateRequest)

		if !ok {
			return nil, fmt.Errorf("invalid request type")
		}

		if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Phone == "" {
			return nil, fmt.Errorf("all fields are required")
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			return nil, err
		}

		return user, nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request type")
		}

		user, err := s.Get(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return user, nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		v, ok := request.(GetAllRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request type")
		}

		filters := Filters{
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Email:     v.Email,
			Phone:     v.Phone,
		}

		// 🎯 Extraemos limit y page directamente del struct GetAllRequest
		// 💡 Si los valores son 0 (no proporcionados), usaremos valores por defecto
		limit := v.Limit
		page := v.Page

		// 🔧 Validación: si limit es 0, usamos el valor por defecto de la configuración
		if limit <= 0 {
			defaultLimit, err := strconv.Atoi(config.LimPageDef)
			if err != nil {
				return nil, fmt.Errorf("invalid default limit configuration: %w", err)
			}
			limit = defaultLimit
		}

		// 🔧 Validación: si page es 0 o negativo, establecemos página 1
		if page <= 0 {
			page = 1
		}

		count, err := s.Count(ctx, filters)

		if err != nil {
			return nil, err
		}

		metaData, err := meta.New(page, limit, int(count), config.LimPageDef)
		if err != nil {
			return nil, err
		}

		users, err := s.GetAll(ctx, filters, metaData.Offset(), metaData.Limit())
		if err != nil {
			return nil, err
		}

		return users, nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(UpdateRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request type")
		}

		// 🎯 Validación: el ID es requerido
		if req.ID == "" {
			return nil, fmt.Errorf("id is required")
		}

		// 🎯 Validación: al menos un campo debe ser proporcionado para actualizar
		if req.FirstName == nil && req.LastName == nil && req.Email == nil && req.Phone == nil {
			return nil, fmt.Errorf("at least one field is required")
		}

		// 🔧 Validación: si se proporciona un campo, no puede estar vacío
		if req.FirstName != nil && *req.FirstName == "" {
			return nil, fmt.Errorf("first name cannot be empty")
		}
		if req.LastName != nil && *req.LastName == "" {
			return nil, fmt.Errorf("last name cannot be empty")
		}
		if req.Email != nil && *req.Email == "" {
			return nil, fmt.Errorf("email cannot be empty")
		}
		if req.Phone != nil && *req.Phone == "" {
			return nil, fmt.Errorf("phone cannot be empty")
		}

		// 💡 Llamamos al servicio para actualizar el usuario
		err = s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			return nil, err
		}

		// ✅ Retornamos un mensaje de éxito (similar al patrón del código antiguo)
		return map[string]string{"message": "User updated successfully"}, nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(DeleteRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request type")
		}

		// 🎯 Validación: el ID es requerido
		if req.ID == "" {
			return nil, fmt.Errorf("id is required")
		}

		// 💡 Llamamos al servicio para eliminar el usuario
		err = s.Delete(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		// ✅ Retornamos un mensaje de éxito
		return map[string]string{"message": "User deleted successfully"}, nil
	}
}
