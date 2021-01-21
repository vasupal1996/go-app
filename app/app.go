package app

import (
	"go-app/server/config"
	mongostorage "go-app/server/storage/mongodb"

	"github.com/rs/zerolog"
)

// Options contains arguments required to create a new app instance
type Options struct {
	MongoDB *mongostorage.MongoStorage
	Logger  *zerolog.Logger
	Config  *config.APPConfig
}

// App := contains resources to implement business logic
type App struct {
	MongoDB *mongostorage.MongoStorage
	Logger  *zerolog.Logger
	Config  *config.APPConfig

	// List of services this app is implementing
	Example Example
}

// NewApp returns new app instance
func NewApp(opts *Options) *App {
	return &App{
		MongoDB: opts.MongoDB,
		Logger:  opts.Logger,
		Config:  opts.Config,
	}
}
