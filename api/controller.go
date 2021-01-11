package api

import (
	"go-app/server/handler"
	"net/http"
)

func (a *API) home(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	requestCTX.SetAppResponse(true, http.StatusAccepted)
	return
}
