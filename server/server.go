/*
	The server package binds all the struct and interfaces of various aspects such as router, database, logging etc.
	StartServer and StopServer functions are exposed to call them via main package or via command line to start/stop
	the execution.

	The server itself listing on some address and port (localhost:8000 (default)) via go routine and will be blocked until
	StopServer function is called via some function or command line. Before stoping server all the resources and connections are closed.
*/

package server

import (
	"context"
	"fmt"
	"go-app/api"
	"go-app/app"
	"go-app/server/auth"
	"go-app/server/config"
	goKafka "go-app/server/kafka"
	"go-app/server/logger"
	"go-app/server/middleware"
	"go-app/server/storage"
	memorystorage "go-app/server/storage/memory"
	mongostorage "go-app/server/storage/mongodb"
	redisstorage "go-app/server/storage/redis"
	"go-app/server/validator"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/urfave/negroni"
)

// Server object encapsulates api, business logic (app),router, storage layer and loggers
type Server struct {
	httpServer *http.Server
	Router     *mux.Router
	Log        *zerolog.Logger
	Config     *config.Config
	Kafka      goKafka.Kafka
	MongoDB    storage.DB
	Redis      storage.Redis

	API *api.API
}

// NewServer returns a new Server object
func NewServer() *Server {
	c := config.GetConfig()
	ms := mongostorage.NewMongoStorage(&c.DatabaseConfig)
	r := mux.NewRouter()

	server := &Server{
		httpServer: &http.Server{},
		Config:     c,
		MongoDB:    ms,
		Router:     r,
	}

	server.InitLoggers()

	if c.ServerConfig.UseMemoryStore {
		server.Redis = memorystorage.NewMemoryStorage()
	} else {
		server.Redis = redisstorage.NewRedisStorage(&c.RedisConfig)
	}

	// Initializing api endpoints and controller
	server.API = api.NewAPI(&api.Options{
		MainRouter: r,
		Logger:     server.Log,
		Config:     &c.APIConfig,
		TokenAuth:  auth.NewTokenAuthentication(&c.TokenAuthConfig),
		Validator:  validator.NewValidation(),
	})

	// Initializing app and services
	server.API.App = app.NewApp(&app.Options{MongoDB: ms, Logger: server.Log, Config: &c.APPConfig})
	// server.API.App.Example = app.InitExample(&app.ExampleOpts{DBName: "example", MongoStorage: ms, Logger: server.Log})

	return server
}

// StartServer function initialize middlewares, router, loggers, and server config. It establishes connection with database & other services.
// After setting up the server it runs the server on specified address and port.
func (s *Server) StartServer() {
	n := negroni.New()

	if s.Config.MiddlewareConfig.EnableRequestLog {
		n.UseFunc(middleware.NewRequestLoggerMiddleware(s.Log).GetMiddlewareHandler())
	}

	n.UseHandler(s.Router)

	s.httpServer = &http.Server{
		Handler:      n,
		Addr:         fmt.Sprintf("%s:%s", s.Config.ServerConfig.ListenAddr, s.Config.ServerConfig.Port),
		ReadTimeout:  s.Config.ServerConfig.ReadTimeout * time.Second,
		WriteTimeout: s.Config.ServerConfig.WriteTimeout * time.Second,
	}

	s.Log.Info().Msgf("Staring server at %s:%s", s.Config.ServerConfig.ListenAddr, s.Config.ServerConfig.Port)
	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil {
			s.Log.Error().Err(err).Msg("")
			return
		}
	}()
}

// StopServer closes all the connection and shutdown the server
func (s *Server) StopServer() {
	if s.Kafka != nil {
		s.Kafka.Close()
	}
	if s.MongoDB != nil {
		s.MongoDB.Close()
	}
	if s.Redis != nil {
		s.Redis.Close()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx)
	os.Exit(0)
}

// InitLoggers initializes all the loggers
func (s *Server) InitLoggers() {
	var kl *logger.KafkaLogWriter
	var cw, fw io.Writer
	if s.Config.LoggerConfig.EnableKafkaLogger {
		dialer := goKafka.NewSegmentioKafkaDialer(&s.Config.KafkaConfig)
		kl = logger.NewKafkaLogWriter(s.Config.LoggerConfig.KafkaLoggerConfig.KafkaTopic, dialer, &s.Config.KafkaConfig)
	}
	if s.Config.LoggerConfig.EnableFileLogger {
		fw = logger.NewFileWriter(s.Config.LoggerConfig.FileLoggerConfig.FileName, s.Config.LoggerConfig.FileLoggerConfig.Path, &s.Config.LoggerConfig.FileLoggerConfig)
	}
	if s.Config.LoggerConfig.EnableConsoleLogger {
		cw = logger.NewZeroLogConsoleWriter(logger.NewStandardConsoleWriter())
	}
	l := logger.NewLogger(kl, cw, fw)

	// Setting logger
	s.Log = l
}
