package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/NicoJCastro/go_lib_response/response"
	"github.com/NicoJCastro/gocourse_user/internal/user"

	"github.com/gorilla/mux"
)

func NewUserHTTPServer(ctx context.Context, endpoints user.Endpoint) http.Handler {
	mux := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	// 🎯 POST /users - Crear usuario
	mux.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeStoreUser,
		encodeResponse,
		opts...,
	)).Methods("POST")

	// 🎯 GET /users/{id} - Obtener un usuario por ID
	mux.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetUser,
		encodeResponse,
		opts...,
	)).Methods("GET")

	// 🎯 GET /users - Obtener todos los usuarios (con paginación y filtros)
	mux.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllUsers,
		encodeResponse,
		opts...,
	)).Methods("GET")

	// 🎯 PATCH /users/{id} - Actualizar usuario
	mux.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateUser,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	// 🎯 DELETE /users/{id} - Eliminar usuario
	mux.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteUser,
		encodeResponse,
		opts...,
	)).Methods("DELETE")
	return mux
}

// 🎯 Decoder para CREATE: decodifica el body JSON
func decodeStoreUser(_ context.Context, r *http.Request) (interface{}, error) {
	var req user.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

// 🎯 Decoder para GET: extrae el ID de la URL
func decodeGetUser(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || id == "" {
		return nil, user.ErrIDRequired
	}
	return user.GetRequest{ID: id}, nil
}

// 🎯 Decoder para GET ALL: extrae query parameters (limit, page, filters)
func decodeGetAllUsers(_ context.Context, r *http.Request) (interface{}, error) {
	// Extraer query parameters
	query := r.URL.Query()

	// Convertir limit y page a int
	limit := 0
	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	page := 0
	if pageStr := query.Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	// Construir GetAllRequest con los query parameters
	req := user.GetAllRequest{
		FirstName: query.Get("first_name"),
		LastName:  query.Get("last_name"),
		Email:     query.Get("email"),
		Phone:     query.Get("phone"),
		Limit:     limit,
		Page:      page,
	}

	return req, nil
}

// 🎯 Decoder para UPDATE: extrae ID de la URL y body JSON
func decodeUpdateUser(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || id == "" {
		return nil, user.ErrIDRequired
	}

	var req user.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, user.ErrInvalidRequestType
	}

	// Asignar el ID extraído de la URL
	req.ID = id
	return req, nil
}

// 🎯 Decoder para DELETE: extrae el ID de la URL
func decodeDeleteUser(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || id == "" {
		return nil, user.ErrIDRequired
	}
	return user.DeleteRequest{ID: id}, nil
}

// 🎯 Encoder para todas las respuestas
func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	respObj := resp.(response.Response)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(respObj.StatusCode())
	return json.NewEncoder(w).Encode(respObj)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 🔍 Intentamos convertir el error a response.Response
	resp, ok := err.(response.Response)
	if !ok {
		// ❌ Si no es response.Response, es un error estándar de Go
		// 💡 Lo convertimos a InternalServerError como fallback seguro
		resp = response.InternalServerError(err.Error())
	}

	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
