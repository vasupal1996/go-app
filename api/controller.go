package api

import (
	"fmt"
	"go-app/app"
	"go-app/server/handler"
	"net/http"
)

func (a *API) home(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	resp, err := a.App.Example.SayHello(false)
	a.Logger.Log().Msg("Here")
	if err != nil {
		requestCTX.SetErr(err, 400)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
	return
}

func (a *API) saveHello(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	opts := app.SaveHelloOpts{}
	err := a.DecodeJSONBody(r, &opts)
	fmt.Println(opts)
	if err != nil {
		fmt.Println(err)
		requestCTX.SetErr(err, 400)
		return
	}
	resp, err := a.App.Example.SaveHello(&opts)
	if err != nil {
		requestCTX.SetErr(err, 400)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusCreated)
	return
}
