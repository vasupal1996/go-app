package handler

import (
	"encoding/json"
	"go-app/server/auth"
	"go-app/server/middleware"
	"net/http"

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
	requestCTX.RequestID = middleware.RequestIDFromContext(r.Context())
	requestCTX.Path = r.URL.Path

	authToken := r.Header.Get("Authorization")
	if authToken != "" {
		err := rh.AuthFunc.VerifyToken(authToken)
		if err != nil {
			requestCTX.SetErr(errors.New("failed to verify token", &errors.PermissionDenied), http.StatusUnauthorized)
			goto SKIP_REQUEST
		} else {
			requestCTX.UserClaim = rh.AuthFunc.GetClaim().(*auth.UserClaim)
		}
	}

	if rh.IsLoggedIn {
		if requestCTX.UserClaim == nil {
			requestCTX.SetErr(errors.New("auth token required", &errors.PermissionDenied), http.StatusUnauthorized)
			goto SKIP_REQUEST
		} else {
			if rh.IsSudoUser {
				if requestCTX.UserClaim.IsAdmin() {
					requestCTX.SetErr(errors.New("permission denied: required admin user role", &errors.PermissionDenied), http.StatusForbidden)
					goto SKIP_REQUEST
				}
			}
		}
	}

SKIP_REQUEST:

	w.Header().Set(auth.HeaderRequestID, requestCTX.RequestID)
	if requestCTX.Err == nil {
		rh.HandlerFunc(requestCTX, w, r)
	}

	if requestCTX.ResponseCode != 0 && requestCTX.ResponseType != RedirectResp {
		w.WriteHeader(requestCTX.ResponseCode)
	}

	switch t := requestCTX.ResponseType; t {
	case HTMLResp:
		w.Header().Set("Content-Type", "text/html")
		res := requestCTX.Response.GetRaw()
		w.Write(res.([]byte))
	case JSONResp:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(requestCTX.Response)
	case ErrorResp:
		w.Header().Set("Content-Type", "application/json")
		requestCTX.Err.RequestID = &requestCTX.RequestID
		json.NewEncoder(w).Encode(&requestCTX.Err)
	case RedirectResp:
		res, _ := requestCTX.Response.MarshalJSON()
		http.Redirect(w, r, string(res), requestCTX.ResponseCode)
	}

}
