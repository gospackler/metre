package logging

import (
	"log"

	"go.uber.org/zap"
)

// Logger is the shared zap logger used within metre.
var Logger *zap.Logger

// Configure the zap logger here to be used within the package.
func init() {
	var err error

	Logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("error generating development logger for metre: %s\n", err)
	}
}
