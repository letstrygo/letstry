package logging

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/letstrygo/letstry/internal/storage"
)

type LogMode int

const (
	LogModeConsole LogMode = iota
	LogModeFile
	LogModeBoth
	LogModeNone
)

func (l LogMode) String() string {
	return [...]string{"console", "file", "both", "none"}[l]
}

type loggingCtxKey string

const (
	loggerCtxKey loggingCtxKey = ".letstry_logger"
	logPrefix    string        = "letstry: "
)

type LoggerConfig struct {
	LogMode LogMode
	Prefix  string
}

type logger struct {
	*log.Logger
	cfg LoggerConfig
}

func (l *logger) Write(p []byte) (n int, err error) {
	l.Logger.Print(string(p))
	return len(p), nil
}

func (l *logger) File() *os.File {
	if file, ok := l.Writer().(*os.File); ok {
		return file
	}

	return nil
}

func (l *logger) ChildLogger(prefix string) *logger {
	return &logger{
		Logger: log.New(l.Writer(), fmt.Sprintf("%s%s: ", logPrefix, prefix), log.LstdFlags),
		cfg:    l.cfg,
	}
}

func New(cfg *LoggerConfig) (*logger, error) {
	storageManager := storage.GetStorage()

	var err error
	var file *os.File
	var internalLogger *log.Logger

	if cfg == nil || cfg.LogMode == LogModeConsole {
		// Write output to the console only.
		internalLogger = log.New(os.Stdout, "letstry: ", log.LstdFlags)
	} else {
		switch cfg.LogMode {
		case LogModeFile:
			file, err = storageManager.OpenFile("ltlog.log")
			internalLogger = log.New(file, fmt.Sprintf("%s%s", logPrefix, cfg.Prefix), log.LstdFlags)
		case LogModeBoth:
			file, err = storageManager.OpenFile("ltlog.log")
			internalLogger = log.New(io.MultiWriter(file, os.Stdout), fmt.Sprintf("%s%s", logPrefix, cfg.Prefix), log.LstdFlags)
		case LogModeNone:
			internalLogger = log.New(io.Discard, "", log.LstdFlags)
		}
	}

	if err != nil {
		return nil, err
	}

	return &logger{
		Logger: internalLogger,
		cfg:    *cfg,
	}, nil
}

func ContextWithLogger(ctx context.Context, logger *logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func LoggerFromContext(ctx context.Context) (*logger, error) {
	logger, ok := ctx.Value(loggerCtxKey).(*logger)
	if !ok {
		return nil, fmt.Errorf("logger not found in context")
	}

	return logger, nil
}

func CloseLog(ctx context.Context) error {
	logger, err := LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	if logFile := logger.File(); logFile != nil {
		return logFile.Close()
	}
	return nil
}
