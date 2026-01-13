package logging

import "go.uber.org/zap"

type Logger struct {
	*zap.Logger
}

func New() (*Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{Logger: logger}, nil
}
