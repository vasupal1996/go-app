package app

import (
	"context"
	"go-app/server/config"
	"go-app/server/logger"
	mongostorage "go-app/server/storage/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewTestApp returns app instance for testing
func NewTestApp(c *config.Config) *App {
	m := mongostorage.NewMongoStorage(&c.DatabaseConfig)
	l := logger.NewLogger(nil, logger.NewZeroLogConsoleWriter(logger.NewStandardConsoleWriter()), nil)
	a := &App{
		MongoDB: m,
		Logger:  l,
		Config:  &c.APPConfig,
	}
	// Setting up services for test app
	// a.Example = InitExample(&ExampleOpts{App: a, DB: m.Client.Database(a.Config.ExampleConfig.DBName), Logger: l})
	return a
}

// CleanTestApp drops the test database
func CleanTestApp(a *App) {
	ctx := context.Background()
	dbs, _ := a.MongoDB.Client.ListDatabases(ctx, bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: "^test_*", Options: "i"}}})
	for _, db := range dbs.Databases {
		a.MongoDB.Client.Database(db.Name).Drop(ctx)
	}
}
