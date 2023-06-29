package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
)

var (
	logger log
)

type log struct {
	zapper *zap.Logger
	path   string
}

type Option = func(*log)

func SetLoggerPath(path string) Option {
	return func(l *log) {
		l.path = path
	}
}

func InitializeLogger(opt ...Option) error {
	for _, option := range opt {
		option(&logger)
	}

	conf := zap.NewProductionConfig()
	if logger.path != "" {
		dir, err := createLogDirectory(logger.path)
		if err != nil {
			return err
		}
		conf.OutputPaths = []string{dir + time.Now().UTC().Format("2006-01-02") + ".log"}
	}

	logg, err := conf.Build()
	if err != nil {
		return err
	}

	logger.zapper = logg

	return nil
}

func createLogDirectory(path string) (string, error) {
	err := os.Mkdir(path, 0700)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	return path + "/", nil
}

func Info(args ...interface{}) {
	logger.zapper.Sugar().Info(args)
}

func Infof(template string, args ...interface{}) {
	logger.zapper.Sugar().Infof(template, args)
}

func Warn(args ...interface{}) {
	logger.zapper.Sugar().Warn(args)
}

func Fatal(args ...interface{}) {
	logger.zapper.Sugar().Fatal(args)
}
