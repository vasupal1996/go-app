package app

// InitService this initializes all the busines logic services
func InitService(a *App) {
	a.Example = InitExample(&ExampleOpts{
		DB:     a.MongoDB.Client.Database(a.Config.ExampleConfig.DBName),
		Logger: a.Logger,
	})
}
