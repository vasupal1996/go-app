package handler

import (
	"encoding/json"

	errors "github.com/vasupal1996/goerror"
)

// ResponseType is the type of response
type ResponseType string

// const are response types constants
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

// AppResponse := app response struct
type AppResponse struct {
	Payload interface{}
}

// MarshalJSON := marshalling error
func (err *AppErr) MarshalJSON() ([]byte, error) {
	var errJSONArr []map[string]interface{}
	for _, err := range err.Error {
		errJSON := errors.Map(err)
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
