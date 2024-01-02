package logging

// Config is the configuration for the logging.
type Config struct {
	// appName is the name of the application.
	appName AppName
}

// NewConfig creates a new Config.
func NewConfig(appName AppName) *Config {
	return &Config{
		appName: appName,
	}
}
