package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NicoJCastro/gocourse_meta/meta"

	"github.com/gorilla/mux"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

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

	UpdateRequest struct {
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
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusBadRequest, Err: "Invalid request payload"})
			return
		}

		if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Phone == "" {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusBadRequest, Err: "All fields are required"})
			return
		}

		user, err := s.Create(req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Err: "Failed to create user"})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: http.StatusCreated, Data: user})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		user, err := s.Get(id)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Err: "Failed to get user"})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: http.StatusOK, Data: user})
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()
		filters := Filters{
			FirstName: v.Get("first_name"),
			LastName:  v.Get("last_name"),
			Email:     v.Get("email"),
			Phone:     v.Get("phone"),
		}

		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		count, err := s.Count(filters)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Err: "Failed to count users"})
			return
		}

		metaData, err := meta.New(page, limit, int(count), config.LimPageDef)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Err: "Failed to create metadata"})
			return
		}

		users, err := s.GetAll(filters, metaData.Offset(), metaData.Limit())
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Err: "Failed to get users"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: http.StatusOK, Data: users, Meta: metaData})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusBadRequest, Err: "Invalid request payload"})
			return
		}

		if req.FirstName != nil && *req.FirstName == "" {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusBadRequest, Err: "First name cannot be empty"})
			return
		}
		if req.LastName != nil && *req.LastName == "" {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusBadRequest, Err: "Last name cannot be empty"})
			return
		}
		path := mux.Vars(r)
		id := path["id"]
		err := s.Update(id, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Err: "Failed to update user"})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: http.StatusOK, Data: map[string]string{"data": "User updated successfully"}})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		err := s.Delete(id)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Err: "Failed to delete user"})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: http.StatusOK, Data: map[string]string{"data": "User deleted successfully"}})
	}
}
