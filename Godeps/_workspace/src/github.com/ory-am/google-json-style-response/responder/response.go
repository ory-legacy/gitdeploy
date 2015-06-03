package responder

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
)

type DataResponse struct {
	ApiVersion string      `json:"apiVersion"`
	Id         uuid.UUID   `json:"id,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

type ErrorItem struct {
	Domain       string `json:"domain,omitempty"`
	Reason       string `json:"reason,omitempty"`
	Message      string `json:"message,omitempty"`
	Location     string `json:"location,omitempty"`
	LocationType string `json:"locationType,omitempty"`
	ExtendedHelp string `json:"extendedHelp,omitempty"`
	SendReport   string `json:"sendReport,omitempty"`
}

type Error struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  []ErrorItem `json:"errors,omitempty"`
}

type ErrorResponse struct {
	ApiVersion string    `json:"apiVersion"`
	Id         uuid.UUID `json:"id,omitempty"`
	Error      Error     `json:"error,omitempty"`
}

type Response interface {
	Marshal() ([]byte, error)
}

func (e ErrorResponse) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e DataResponse) Marshal() ([]byte, error) {
	return json.Marshal(e)
}
