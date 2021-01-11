package handler

import (
	"encoding/json"
	"reflect"

	errors "github.com/vasupal1996/goerror"
)

// ResponseType custom types for defining response type
type ResponseType string

// responses type values for multiple responses
const (
	HTMLResp     ResponseType = "html"
	JSONResp     ResponseType = "json"
	RedirectResp ResponseType = "redirect"
	FileResp     ResponseType = "file"
	ErrorResp    ResponseType = "error"
)

// AppErr := app error struct
type AppErr struct {
	Error     []error
	RequestID *string
}

// Response := response interface
type Response interface {
	MarshalJSON() ([]byte, error)
	GetRaw() interface{}
}

// AppResponse := app response struct
type AppResponse struct {
	Payload interface{}
}

// CustomResponse := app response struct with success
type CustomResponse struct {
	Success bool
	Payload interface{}
	Error   interface{}
}

// MarshalJSON := marshalling error
func (err *AppErr) MarshalJSON() ([]byte, error) {
	var errJSONArr []map[string]interface{}
	for _, e := range err.Error {
		errJSON := errors.Map(e)
		errJSONArr = append(errJSONArr, errJSON)
	}
	return json.Marshal(&struct {
		Error     []map[string]interface{} `json:"error"`
		Success   bool                     `json:"success"`
		RequestID *string                  `json:"request_id"`
	}{
		Error:     errJSONArr,
		Success:   false,
		RequestID: err.RequestID,
	})
}

// MarshalJSON := marshalling error
func (r *AppResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Success bool        `json:"success"`
		Payload interface{} `json:"payload"`
	}{
		Success: true,
		Payload: &r.Payload,
	})
}

// GetRaw returns payload
func (r *AppResponse) GetRaw() interface{} {
	return r.Payload
}

// MarshalJSON := marshalling error
func (r *CustomResponse) MarshalJSON() ([]byte, error) {
	var errJSON []map[string]interface{}
	if reflect.TypeOf(r.Error).Kind() == reflect.Slice || reflect.TypeOf(r.Error).Kind() == reflect.Array {
		for _, e := range r.Error.([]error) {
			errJSON = append(errJSON, errors.Map(e))
		}
	} else {
		errJSON = append(errJSON, errors.Map(r.Error.(error)))
	}

	return json.Marshal(&struct {
		Success bool        `json:"success"`
		Payload interface{} `json:"payload"`
		Error   interface{} `json:"error"`
	}{
		Success: r.Success,
		Payload: &r.Payload,
		Error:   errJSON,
	})
}

// GetRaw returns payload
func (r *CustomResponse) GetRaw() interface{} {
	return r.Payload
}
