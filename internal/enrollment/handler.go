package enrollment

import (
	"encoding/json"
	"errors"
	"goWeb/pkg/meta"
	"net/http"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoint struct {
		Create Controller
	}

	CreateRequest struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
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
		Create: makeCreateHandler(s),
	}
}

func makeCreateHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Status: http.StatusBadRequest,
				Error:  "invalid request payload",
			})
			return
		}

		if req.UserID == "" || req.CourseID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Status: http.StatusBadRequest,
				Error:  "user_id and course_id are required",
			})
			return
		}

		enrollment, err := s.Create(req.UserID, req.CourseID)
		if err != nil {
			statusCode := http.StatusInternalServerError
			errorMsg := err.Error()

			if errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrCourseNotFound) {
				statusCode = http.StatusNotFound
			}

			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(Response{
				Status: statusCode,
				Error:  errorMsg,
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&Response{
			Status: http.StatusCreated,
			Data:   enrollment,
		})
	}
}
