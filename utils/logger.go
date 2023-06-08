package utils

import (
	"log"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitializeLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	Logger = logger
}
