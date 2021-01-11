package handler

import (
	"encoding/json"
	"go-app/server/auth"
	"net/http"

	uuid "github.com/satori/go.uuid"
	errors "github.com/vasupal1996/goerror"
)

// Request represents a request from client
type Request struct {
	HandlerFunc func(*RequestContext, http.ResponseWriter, *http.Request)
	AuthFunc    auth.TokenAuth
	IsLoggedIn  bool
	IsSudoUser  bool
}

// HandleRequest := handles incoming requests from client
func (rh *Request) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestCTX := &RequestContext{}
	requestCTX.RequestID = uuid.NewV4().String() + "-" + uuid.NewV4().String()
	requestCTX.Path = r.URL.Path

	authToken := r.Header.Get("Authorization")
	if authToken != "" {
		err := rh.AuthFunc.VerifyToken(authToken)
		if err != nil {
			requestCTX.SetErr(errors.New("failed to verify token", &errors.PermissionDenied))
			goto SKIP_REQUEST
		} else {
			requestCTX.UserClaim = rh.AuthFunc.GetClaim().(*auth.UserClaim)
		}
	}

	if rh.IsLoggedIn {
		if requestCTX.UserClaim == nil {
			requestCTX.SetErr(errors.New("auth token required", &errors.PermissionDenied))
			goto SKIP_REQUEST
		} else {
			if rh.IsSudoUser {
				if requestCTX.UserClaim.IsAdmin() {
					requestCTX.SetErr(errors.New("permission denied: required admin user role", &errors.PermissionDenied))
					goto SKIP_REQUEST
				}
			}
		}
	}

SKIP_REQUEST:

	if requestCTX.Err == nil {
		w.Header().Set(auth.HeaderRequestID, requestCTX.RequestID)
		rh.HandlerFunc(requestCTX, w, r)
	}

	if requestCTX.ResponseCode != 0 && requestCTX.ResponseType != RedirectResp {
		w.WriteHeader(requestCTX.ResponseCode)
	}

	switch t := requestCTX.ResponseType; t {
	case HTMLResp:
		w.Header().Set("Content-Type", "text/html")
		w.Write(requestCTX.Response.Payload.([]byte))
	case JSONResp:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(requestCTX.Response)
	case ErrorResp:
		w.Header().Set("Content-Type", "application/json")
		requestCTX.Err.RequestID = &requestCTX.RequestID
		json.NewEncoder(w).Encode(&requestCTX.Err)
	case RedirectResp:
		http.Redirect(w, r, requestCTX.Response.Payload.(string), requestCTX.ResponseCode)
	}
}
