package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
)

var (
	Logger logger
)

type logger struct {
	logger *zap.Logger
	path   string
}

type Option = func(*logger)

func SetLoggerPath(path string) Option {
	return func(l *logger) {
		l.path = path
	}
}

func InitializeLogger(opt ...Option) error {
	for _, option := range opt {
		option(&Logger)
	}

	conf := zap.NewProductionConfig()
	if Logger.path != "" {
		dir, err := createLogDirectory(Logger.path)
		if err != nil {
			return err
		}
		conf.OutputPaths = []string{dir + time.Now().UTC().Format("2006-01-02") + ".log"}
	}

	logg, err := conf.Build()
	if err != nil {
		return err
	}

	Logger.logger = logg

	return nil
}

func createLogDirectory(path string) (string, error) {
	err := os.Mkdir(path, 0700)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	return path + "/", nil
}

func (l *logger) Info(args ...interface{}) {
	l.logger.Sugar().Info(args)
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.logger.Sugar().Infof(template, args)
}

func (l *logger) Warn(args ...interface{}) {
	l.logger.Sugar().Warn(args)
}

func (l *logger) Fatal(args ...interface{}) {
	l.logger.Sugar().Fatal(args)
}
