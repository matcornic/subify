package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/vincentdaniel/subify/common/config"
	"os"
	logger "github.com/spf13/jwalterweatherman"
	"log"
)

//GetHashOfVideo gets the hash used by SubDb to identify a video. Absolutely needed either to download or upload subtitles.
//The hash is composed by taking the first and the last 64kb of the video file, putting all together and generating a md5 of the resulting data (128kb).
func GetHashOfVideo(filename string) string {
	readsize := 64 * 1024 // 64kb

	// Open Video
	file, err := os.Open(filename)
	if err != nil {
		ExitPrintError(err, "Can't open file ", filename)
	}
	defer file.Close()

	// Get stats of file
	fi, err := file.Stat()
	if err != nil {
		ExitPrintError(err, "Can't get stats for file ", filename)
	}

	// Fill a buffer with first bytes of file
	bufB := make([]byte, readsize)
	_, err = file.Read(bufB)
	if err != nil {
		ExitPrintError(err, "Can't read content of file ", filename)
	}

	//Fill a buffer with last bytes of file
	bufE := make([]byte, readsize)
	n, err := file.ReadAt(bufE, fi.Size()-int64(len(bufE)))
	if err != nil {
		ExitPrintError(err, "Can't read content of file. File is probably too small: ", filename)
	}
	bufE = bufE[:n]

	// Generates MD5 of both bytes chain
	bufB = append(bufB, bufE...)
	hash := fmt.Sprintf("%x", md5.Sum(bufB))

	VerbosePrintln(logger.INFO, "Hash of video is " + hash)

	return hash
}

func VerbosePrintln(logger *log.Logger, log string) {
	if config.Verbose {
		logger.Println(log)
	}
}

func Exit(format string, args ...interface{}) {
	ExitVerbose("", format, args...)
}

/*
 * Exit func displays an error message on stderr and exit 1
 * Eventually prints more details about the error if verbose mode is enabled
 */
func ExitPrintError(err error, format string, args ...interface{}) {
	ExitVerbose(fmt.Sprint(err), format, args...)
}


/*
 * Exit func displays an error message on stderr and exit 1
 * Eventually prints more details if any verbose details are given and verbose mode is enabled
 */
func ExitVerbose(verbose string, format string, args ...interface{}) {
	if verbose != "" {
		VerbosePrintln(logger.ERROR, verbose)
	} else if !config.Verbose {
		logger.ERROR.Println("Run subify with --verbose option to get more information about the error")
	}
	logger.FATAL.Printf(format)
}
