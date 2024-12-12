package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
}

func New(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)

	return logger
}

func WithFile(logger *logrus.Logger, fd *os.File) *logrus.Logger {
	logger.SetOutput(fd)

	return logger
}

func WithFormat(logger *logrus.Logger, format *logrus.JSONFormatter) *logrus.Logger {
	logger.SetFormatter(format)

	return logger
}

func CreateAndOpen(path string) (*os.File, error) {
	dirs := strings.Split(path, "/")
	dirpath := strings.Join(dirs[:len(dirs)-1], "/")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
			return nil, fmt.Errorf("couldn't create a log dir: %v", err)
		}
	}

	fd, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return fd, nil
}
