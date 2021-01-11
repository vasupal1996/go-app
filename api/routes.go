package api

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {
	a.Router.Root.Handle("/", a.requestHandler(a.home)).Methods("GET")
}

// InitTestRoutes := intializing all the testing and development endpoints
func (a *API) InitTestRoutes() {
	// a.Router.APIRoot.Handle("/test/add-category", a.requestHandler(a.addSampleCategories)).Methods("GET")
	// a.Router.APIRoot.Handle("/test/add-catalog-and-variants", a.requestHandler(a.imageUpload)).Methods("GET")
}

// InitMediaRoutes initialize media urls
func (a *API) InitMediaRoutes() {
	// fs := http.FileServer(http.Dir("./.media/"))
	// a.Router.Root.PathPrefix("/media/").Handler(http.StripPrefix("/media/", fs))
}
