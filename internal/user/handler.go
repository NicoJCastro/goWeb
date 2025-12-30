package user

import (
	"encoding/json"
	"net/http"

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

	ErrorResponse struct {
		Error string `json:"error"`
	}
)

func MakeEndpoints(s Service) Endpoint {
	return Endpoint{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request payload"})
			return
		}

		if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Phone == "" {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "All fields are required"})
			return
		}

		user, err := s.Create(req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create user"})
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		user, err := s.Get(id)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to get user"})
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.GetAll()
		if err != nil {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to get users"})
			return
		}
		json.NewEncoder(w).Encode(users)
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request payload"})
			return
		}

		if req.FirstName != nil && *req.FirstName == "" {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "First name cannot be empty"})
			return
		}
		if req.LastName != nil && *req.LastName == "" {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Last name cannot be empty"})
			return
		}
		path := mux.Vars(r)
		id := path["id"]
		err := s.Update(id, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update user"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"data": "User updated successfully"})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		err := s.Delete(id)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to delete user"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"data": "User deleted successfully"})
	}
}
