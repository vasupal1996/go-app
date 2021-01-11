package handler

import (
	"go-app/server/auth"
)

// RequestContext persists the request journey from request hitting the server to sending response
type RequestContext struct {
	RequestID    string
	Path         string
	Response     Response
	Err          *AppErr
	ResponseType ResponseType
	ResponseCode int
	UserClaim    auth.Claim
}

// SetErr := setting Err response in request context
func (requestCTX *RequestContext) SetErr(err error, statusCode int) {
	appErr := requestCTX.Err
	requestCTX.ResponseType = ErrorResp
	requestCTX.ResponseCode = statusCode
	if appErr == nil {
		appErr = &AppErr{}
	}
	appErr.Error = append(appErr.Error, err)
	requestCTX.Err = appErr
}

// SetAppResponse := setting app response in request context
func (requestCTX *RequestContext) SetAppResponse(message interface{}, statusCode int) {
	requestCTX.ResponseType = JSONResp
	requestCTX.ResponseCode = statusCode
	requestCTX.Response = &AppResponse{
		Payload: message,
	}
}

// SetCustomResponse := setting app response in request context
func (requestCTX *RequestContext) SetCustomResponse(success bool, message interface{}, err interface{}, statusCode int) {
	requestCTX.ResponseType = JSONResp
	requestCTX.ResponseCode = statusCode
	requestCTX.Response = &CustomResponse{
		Success: success,
		Payload: message,
		Error:   err,
	}
}

// SetErrs adds slice of errors in requestCTX
func (requestCTX *RequestContext) SetErrs(errs []error, statusCode int) {
	for _, e := range errs {
		requestCTX.SetErr(e, statusCode)
	}
}

// SetHTMLResponse := setting app html response in request context
func (requestCTX *RequestContext) SetHTMLResponse(message []byte, statusCode int) {
	requestCTX.ResponseType = HTMLResp
	requestCTX.ResponseCode = statusCode
	requestCTX.Response = &AppResponse{
		Payload: message,
	}
}

// SetRedirectResponse := setting app redirect response in request context
func (requestCTX *RequestContext) SetRedirectResponse(message string, statusCode int) {
	requestCTX.ResponseType = RedirectResp
	requestCTX.ResponseCode = statusCode
	requestCTX.Response = &AppResponse{
		Payload: message,
	}
}
