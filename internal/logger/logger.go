package logger

import (
	"fmt"
	"github.com/dnsoftware/gophermart/internal/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"sync"
)

type logger struct {
	filename string
	*zap.Logger
}

var projectLogger *logger
var once sync.Once

// Log Получение синглтона логгера
func Log() *logger {
	once.Do(func() {
		var err error

		projectLogger, err = createLogger(constants.LogFile, constants.LogLevel)
		if err != nil {
			log.Fatal(err)
		}
	})

	return projectLogger
}

// логирование в файл и в консоль
func createLogger(filename string, logLevel zapcore.Level) (*logger, error) {
	// формат времени "2006-01-02T15:04:05.000Z0700"
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	// создание кодировщиков для вывода в файл и в консоль
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Открываем лог файл
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create writers for file and console
	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Create cores for writing to the file and console
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, logLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, logLevel)

	// Combine cores
	core := zapcore.NewTee(fileCore, consoleCore)

	// Create the logger with additional context information (caller, stack trace)
	l := &logger{
		filename,
		zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)),
	}

	return l, nil
}
