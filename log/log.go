package log

import (
	"github.com/uber-go/zap"
)

// Logger is the logger used in the application
var Logger zap.Logger

func init() {
	Logger = zap.New(
		zap.NewJSONEncoder(zap.NoTime()), // drop timestamps in tests
	)
}
