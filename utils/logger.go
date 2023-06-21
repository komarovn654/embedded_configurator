package utils

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
)

func InitializeLogger() {
	conf := zap.NewDevelopmentConfig()
	conf.OutputPaths = []string{createLogDirectory() + time.Now().UTC().Format("2006-01-02") + "_log.txt"}
	conf.Development = true

	logger, err := conf.Build()
	if err != nil {
		log.Fatal(err)
	}
	Logger = logger

}

func createLogDirectory() string {
	path := ".log"
	err := os.Mkdir(path, 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	return path + "/"
}
