package handler

import "go-app/server/auth"

// RequestContext persists the request journey from request hitting the server to sending response
type RequestContext struct {
	RequestID    string
	Path         string
	Response     *AppResponse
	Err          *AppErr
	ResponseType ResponseType
	ResponseCode int
	UserClaim    auth.Claim
}

// SetErr := setting Err response in request context
func (requestCTX *RequestContext) SetErr(err error) {
	appErr := requestCTX.Err
	requestCTX.ResponseType = ErrorResp
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

// SetErrs adds slice of errors in requestCTX
func (requestCTX *RequestContext) SetErrs(errs []error) {
	for _, e := range errs {
		requestCTX.SetErr(e)
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
