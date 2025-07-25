package util

import (
	"os"
	"github.com/rs/zerolog"
)

type CustomLogger struct {
	*zerolog.Logger
	file *os.File
}

func NewCustomFileLogger(debug bool) (*CustomLogger, error) {
	file, err := os.Create("logfile.log")
	if err != nil {
		return nil, err
	}
	
  var logLevel zerolog.Level
  
  if debug == true {
    logLevel = zerolog.DebugLevel
  } else {
    logLevel = zerolog.InfoLevel
  }

  logger := zerolog.New(file).With().Timestamp().Logger().Level(logLevel)
	customLogger := &CustomLogger{
		Logger: &logger,
		file:   file,
	}
	return customLogger, nil
}

func NewCustomLogger(debug bool) (*CustomLogger, error) {
  var logLevel zerolog.Level
  
  if debug == true {
    logLevel = zerolog.DebugLevel
  } else {
    logLevel = zerolog.InfoLevel
  }

  logger := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(logLevel)
	customLogger := &CustomLogger{
		Logger: &logger,
	}
	return customLogger, nil
}


func (cl *CustomLogger) Close() {
	if cl.file != nil {
		cl.file.Close()
	}
}
