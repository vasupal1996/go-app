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
	"go-app/server/config"
	goKafka "go-app/server/kafka"
	"go-app/server/logger"
	"go-app/server/middleware"
	"go-app/server/storage"
	memorystorage "go-app/server/storage/memory"
	mongostorage "go-app/server/storage/mongodb"
	redisstorage "go-app/server/storage/redis"
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
	Kafka      goKafka.KImpl
	MongoDB    storage.Impl
	Redis      storage.Impl

	API *api.API
}

// NewServer returns a new Server object
func NewServer() *Server {
	c := config.GetConfig()
	k := goKafka.NewKafka(&c.KafkaConfig)
	l := logger.NewLogger(&c.LoggerConfig, logger.NewKafkaLogWriter(c.LoggerConfig.KafkaLoggerConfig.KafkaTopic, k.Conn), logger.NewZeroLogConsoleLogger(logger.NewStandardConsoleWriter()), nil)
	ms := mongostorage.NewMongoStorage(&c.DatabaseConfig)
	r := mux.NewRouter()

	api := api.NewAPI(r, &c.APIConfig, &c.TokenAuthConfig, l)

	server := &Server{
		httpServer: &http.Server{},
		Config:     c,
		Kafka:      k,
		Log:        l,
		MongoDB:    ms,
		Router:     r,
		API:        api,
	}

	if c.ServerConfig.UserMemoryStore {
		server.Redis = memorystorage.NewMemoryStorage()
	} else {
		server.Redis = redisstorage.NewRedisStorage(&c.RedisConfig)
	}
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
	s.Kafka.Close()
	s.MongoDB.Close()
	s.Redis.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx)
	os.Exit(0)
}
