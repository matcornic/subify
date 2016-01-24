package utils

import (
	"fmt"
	"log"

	"github.com/matcornic/subify/common/config"
	logger "github.com/spf13/jwalterweatherman"
)

// VerbosePrintln only prints log if verbose mode is enabled
func VerbosePrintln(logger *log.Logger, log string) {
	if config.Verbose && log != "" {
		logger.Println(log)
	}
}

// Exit exits the application and logs the given message
func Exit(format string, args ...interface{}) {
	ExitVerbose("", format, args...)
}

// ExitPrintError displays an error message on stderr and exit 1
// Eventually prints more details about the error if verbose mode is enabled
func ExitPrintError(err error, format string, args ...interface{}) {
	ExitVerbose(fmt.Sprint(err), format, args...)
}

// ExitVerbose displays an error message on stderr and exit 1
// Eventually prints more details if any verbose details are given and verbose mode is enabled
func ExitVerbose(verboseLog string, format string, args ...interface{}) {
	VerbosePrintln(logger.ERROR, verboseLog)
	if !config.Verbose {
		logger.ERROR.Println("Run subify with --verbose option to get more information about the error")
	}
	logger.FATAL.Printf(format)
}
