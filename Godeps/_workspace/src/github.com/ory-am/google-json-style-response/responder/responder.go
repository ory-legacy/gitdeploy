package responder

import (
	"code.google.com/p/go-uuid/uuid"
	"net/http"
)

type responder struct {
	apiVersion    string
	errorResponse *ErrorResponse
	dataResponse  *DataResponse
}

func New(apiVersion string) *responder {
	return &responder{
		apiVersion: apiVersion,
		errorResponse: &ErrorResponse{
			Error: Error{
				Errors: make([]ErrorItem, 0, 1),
			},
		},
		dataResponse: &DataResponse{},
	}
}

func (r *responder) Success(data interface{}) Response {
	r.dataResponse = &DataResponse{
		ApiVersion: r.apiVersion,
		Id:         uuid.NewRandom(),
		Data:       data,
	}
	return *r.dataResponse
}

func (r *responder) AddError(e ErrorItem) {
	r.errorResponse.Error.Errors = append(r.errorResponse.Error.Errors, e)
}

func (r *responder) Error(code int, message string) Response {
	r.errorResponse = &ErrorResponse{
		ApiVersion: r.apiVersion,
		Id:         uuid.NewRandom(),
		Error: Error{
			Code:    code,
			Message: message,
		},
	}
	return r.errorResponse
}

func (r *responder) Write(w http.ResponseWriter, s Response) error {
	m, err := s.Marshal()
	if err != nil {
		responseError := r.Error(http.StatusInternalServerError, err.Error())
		return r.Write(w, responseError)
	}
	w.Header().Set("Content-Type", "application/json")
	switch s := s.(type) {
	case *ErrorResponse:
		if s.Error.Code == 0 {
			s.Error.Code = http.StatusInternalServerError
		}
		code := s.Error.Code
		w.WriteHeader(code)
	}
	w.Write(m)
	return nil
}
