package logging

import (
	"fmt"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"strings"
)

func CommonLogger(cfg *Config) (*slog.Logger, error) {
	var writer io.Writer
	if cfg != nil {
		logFile, err := os.OpenFile(fmt.Sprintf("%s%s-logs.log", cfg.logFilePath, cfg.appName),
			os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
		writer = io.MultiWriter(os.Stdout, logFile)
	} else {
		writer = os.Stdout
	}
	return CommonLoggerWithOptions(cfg, writer, slog.LevelDebug, true)
}

// CommonLoggerWithOptions constructs a logger with specific options.
func CommonLoggerWithOptions(cfg *Config, w io.Writer, level slog.Level, logToJson bool) (*slog.Logger, error) {
	opts := slog.HandlerOptions{
		AddSource:   true,
		Level:       level,
		ReplaceAttr: replaceAttrs,
	}

	var logger *slog.Logger

	if logToJson {
		logger = slog.New(opts.NewJSONHandler(w))
	} else {
		logger = slog.New(opts.NewTextHandler(w))
	}

	logger = logger.With(
		"app", cfg.appName,
	)

	slog.SetDefault(logger)
	return logger, nil
}

// replaceAttrs replaces select attribute fields in an incoming log record to suit logging standards as required.
func replaceAttrs(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.SourceKey:
		// Cut the source file to a relative path.
		v := strings.Split(a.Value.String(), "/")
		idx := len(v) - 2
		if idx < 0 {
			idx = 0
		}
		a.Value = slog.StringValue(strings.Join(v[idx:], "/"))
	}
	return a
}
