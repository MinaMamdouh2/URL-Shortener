package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(service string, outputPaths ...string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	// ISO8601 (e.g., 2025-04-21T15:04:05Z07:00)
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// Disable stack tracing in production to reduce noise unless we need to explicitly
	// want stack traces for debugging
	config.DisableStacktrace = true
	// Adds initial fields to every log entry, here we are adding "service" name
	// to identify which service is using the logger
	config.InitialFields = map[string]any{
		"service": service,
	}

	config.OutputPaths = []string{"stdout"}
	if outputPaths != nil {
		config.OutputPaths = outputPaths
	}

	// Builds the logger from the config, "zap.WithCaller(true)" this includes file
	// & line number info into logs (e.g., '"caller":"main.go:42")
	log, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}
	return log.Sugar(), nil
}
