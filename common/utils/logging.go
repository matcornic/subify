package utils

import (
	"log"
	"os"

	"github.com/matcornic/subify/common/config"
	logger "github.com/spf13/jwalterweatherman"
)

// InitLoggingConf initializes the loggers used in the rest of the application
func InitLoggingConf() {

	errorLoggers := []*log.Logger{logger.ERROR, logger.CRITICAL, logger.FATAL}
	allLoggers := append(errorLoggers, []*log.Logger{logger.INFO, logger.WARN}...)

	// Log levels displayed to the user should not include debug/trace information
	// hence we remove the flags
	for _, log := range allLoggers {
		log.SetFlags(0)
	}

	for _, log := range errorLoggers {
		log.SetOutput(os.Stderr)
	}

	if config.Verbose {
		logger.INFO.SetOutput(os.Stdout)
		logger.WARN.SetOutput(os.Stdout)
	}
}
