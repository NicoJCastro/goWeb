package course

import (
	"encoding/json"
	"goWeb/pkg/meta"
	"net/http"
	"strconv"

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

	CreateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	UpdateReq struct {
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Error  string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
)

func MakeEndpoint(s Service) Endpoint {
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

		var req CreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Error:  "invalid request payload",
			})
			return
		}

		if req.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Error:  "name is required",
			})
			return
		}

		if req.StartDate == "" || req.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Error:  "start_date and end_date are required",
			})
			return
		}

		course, err := s.Create(req.Name, req.StartDate, req.EndDate)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Error:  "failed to create course",
			})
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   course,
		})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		course, err := s.Get(id)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Error: "Failed to get course"})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: http.StatusOK, Data: course})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()
		filters := Filters{
			Name: v.Get("name"),
		}

		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		count, err := s.Count(filters)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Error: "Failed to count users"})
			return
		}

		metaData, err := meta.New(page, limit, int(count))
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Error: "Failed to create metadata"})
			return
		}

		users, err := s.GetAll(filters, metaData.Offset(), metaData.Limit())
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Error: "Failed to get users"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: http.StatusOK, Data: users, Meta: metaData})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusBadRequest, Error: "Invalid request payload"})
			return
		}

		if req.Name == nil && req.StartDate == nil && req.EndDate == nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusBadRequest, Error: "At least one field (name, start_date, end_date) must be provided for update"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]
		err := s.Update(id, req.Name, req.StartDate, req.EndDate)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Error: "Failed to update course"})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: http.StatusOK, Data: map[string]string{"data": "Course updated successfully"}})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		err := s.Delete(id)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: http.StatusInternalServerError, Error: "Failed to delete course"})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: http.StatusOK, Data: map[string]string{"data": "Course deleted successfully"}})
	}
}
