package logging

type LoggerName string
type LoggerPath string

type Config struct {
	appName     LoggerName
	logFilePath LoggerPath
}

func NewConfig(appName LoggerName, logFilePath LoggerPath) *Config {
	return &Config{
		appName:     appName,
		logFilePath: logFilePath,
	}
}
