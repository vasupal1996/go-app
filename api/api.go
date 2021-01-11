package api

import (
	"go-app/server/auth"
	"go-app/server/config"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

// API := returns API struct
type API struct {
	Router     *Router
	MainRouter *mux.Router
	Logger     *zerolog.Logger
	Config     *config.APIConfig
	TokenAuth  auth.TokenAuth
}

// Router stores all the endpoints available for the server to respond.
type Router struct {
	Root       *mux.Router
	APIRoot    *mux.Router
	StaticRoot *mux.Router
}

// NewAPI returns API instance
func NewAPI(m *mux.Router, c *config.APIConfig, tc *config.TokenAuthConfig, l *zerolog.Logger) *API {
	ta := auth.NewTokenAuthentication(tc)
	api := API{
		MainRouter: m,
		Router:     &Router{},
		Config:     c,
		TokenAuth:  ta,
		Logger:     l,
	}
	api.setupRoutes()
	return &api
}

func (a *API) setupRoutes() {
	a.Router.Root = a.MainRouter
	a.Router.APIRoot = a.MainRouter.PathPrefix("/api").Subrouter()
	a.InitRoutes()
	if a.Config.EnableTestRoute {
		a.InitTestRoutes()
	}
	if a.Config.EnableMediaRoute {
		a.InitMediaRoutes()
	}
	if a.Config.EnableStaticRoute {
		a.Router.StaticRoot = a.MainRouter.PathPrefix("/static").Subrouter()
	}
}

func (a *API) requestHandler(h func(c *handler.RequestContext, w http.ResponseWriter, r *http.Request)) http.Handler {
	return &handler.Request{
		HandlerFunc: h,
		AuthFunc:    a.TokenAuth,
		IsLoggedIn:  false,
		IsSudoUser:  false,
	}
}

func (a *API) requestWithAuthHandler(h func(c *handler.RequestContext, w http.ResponseWriter, r *http.Request)) http.Handler {
	return &handler.Request{
		HandlerFunc: h,
		AuthFunc:    a.TokenAuth,
		IsLoggedIn:  true,
		IsSudoUser:  false,
	}
}

func (a *API) requestWithSudoHandler(h func(c *handler.RequestContext, w http.ResponseWriter, r *http.Request)) http.Handler {
	return &handler.Request{
		HandlerFunc: h,
		AuthFunc:    a.TokenAuth,
		IsLoggedIn:  true,
		IsSudoUser:  true,
	}
}