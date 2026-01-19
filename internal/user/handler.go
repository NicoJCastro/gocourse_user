package user

import (
	"context"
	"strconv"

	"github.com/NicoJCastro/go_lib_response/response"
	"github.com/NicoJCastro/gocourse_meta/meta"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateRequest)

		if !ok {
			return nil, response.BadRequest("invalid request type")
		}

		if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Phone == "" {
			return nil, response.BadRequest("all fields are required")
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("User created successfully", user, nil), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetRequest)
		if !ok {
			return nil, response.BadRequest("invalid request type")
		}

		user, err := s.Get(ctx, req.ID)
		if err != nil {
			return nil, response.NotFound(err.Error())
		}

		return response.OK("User retrieved successfully", user, nil), nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		v, ok := request.(GetAllRequest)
		if !ok {
			return nil, response.BadRequest("invalid request type")
		}

		filters := Filters{
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Email:     v.Email,
			Phone:     v.Phone,
		}

		// Extraemos limit y page directamente del struct GetAllRequest
		// Si los valores son 0 (no proporcionados), usaremos valores por defecto
		limit := v.Limit
		page := v.Page

		// 🔧 Validación: si limit es 0, usamos el valor por defecto de la configuración
		if limit <= 0 {
			defaultLimit, err := strconv.Atoi(config.LimPageDef)
			if err != nil {
				return nil, response.InternalServerError(err.Error())
			}
			limit = defaultLimit
		}

		// 🔧 Validación: si page es 0 o negativo, establecemos página 1
		if page <= 0 {
			page = 1
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		metaData, err := meta.New(page, limit, int(count), config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		users, err := s.GetAll(ctx, filters, metaData.Offset(), metaData.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("Users retrieved successfully", users, metaData), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdateRequest)
		if !ok {
			return nil, response.BadRequest("invalid request type")
		}

		// 🎯 Validación: el ID es requerido
		if req.ID == "" {
			return nil, response.BadRequest("id is required")
		}

		// 🎯 Validación: al menos un campo debe ser proporcionado para actualizar
		if req.FirstName == nil && req.LastName == nil && req.Email == nil && req.Phone == nil {
			return nil, response.BadRequest("at least one field is required")
		}

		// 🔧 Validación: si se proporciona un campo, no puede estar vacío
		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest("first name cannot be empty")
		}
		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest("last name cannot be empty")
		}
		if req.Email != nil && *req.Email == "" {
			return nil, response.BadRequest("email cannot be empty")
		}
		if req.Phone != nil && *req.Phone == "" {
			return nil, response.BadRequest("phone cannot be empty")
		}

		// ✅ Llamamos al servicio y obtenemos el usuario actualizado
		user, err := s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		// ✅ Retornamos el usuario actualizado en la respuesta
		return response.OK("User updated successfully", user, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DeleteRequest)
		if !ok {
			return nil, response.BadRequest("invalid request type")
		}

		// 🎯 Validación: el ID es requerido
		if req.ID == "" {
			return nil, response.BadRequest("id is required")
		}

		// 💡 Llamamos al servicio para eliminar el usuario
		err := s.Delete(ctx, req.ID)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		// ✅ Retornamos un mensaje de éxito
		return response.OK("User deleted successfully", nil, nil), nil
	}
}
