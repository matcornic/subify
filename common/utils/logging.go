package utils

import (
	logger "github.com/spf13/jwalterweatherman"
	"os"
	"log"
)

// TODO table with all error loggers
// TODO table with all loggers

func InitLoggingConf() {

	errorLoggers := []*log.Logger{logger.ERROR, logger.CRITICAL, logger.FATAL}
	allLoggers := append(errorLoggers, []*log.Logger{logger.INFO, logger.WARN}...)

	// Log levels displayed to the user should not include debug/trace information, hence we remove the flags
	for _, log := range allLoggers {
		log.SetFlags(0)
	}

	for _, log := range errorLoggers {
		log.SetOutput(os.Stderr)
	}
}

