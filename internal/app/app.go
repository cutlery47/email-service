package app

import (
	"context"
	"fmt"
	"time"

	"github.com/cutlery47/email-service/internal/config"
	v1 "github.com/cutlery47/email-service/internal/controller/http/v1"
	"github.com/cutlery47/email-service/internal/repo"
	"github.com/cutlery47/email-service/internal/service"
	"github.com/cutlery47/email-service/pkg/httpserver"
	"github.com/cutlery47/email-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"
)

func Run() error {
	rand.Seed(uint64(time.Now().UnixNano()))
	ctx := context.Background()

	logrus.Debug("reading config...")
	conf, err := config.New()
	if err != nil {
		return fmt.Errorf("error when setting up config: %v", err)
	}

	logrus.Debug("creating loggers...")
	infoFd, err := logger.CreateAndOpen(conf.Logger.InfoPath)
	if err != nil {
		return fmt.Errorf("error when creating info log file: %v", err)
	}

	errFd, err := logger.CreateAndOpen(conf.Logger.ErrorPath)
	if err != nil {
		return fmt.Errorf("error when creating error log file: %v", err)
	}

	infoLog := logger.WithFormat(logger.WithFile(logger.New(logrus.InfoLevel), infoFd), &logrus.JSONFormatter{})
	errLog := logger.WithFormat(logger.WithFile(logger.New(logrus.ErrorLevel), errFd), &logrus.JSONFormatter{})

	logrus.Debug("initializing cache...")
	cache := service.NewMapCache(conf.Cache)

	logrus.Debug("initializing repository...")
	repo := repo.NewMailRepository(ctx, conf.Postgres)

	logrus.Debug("initializing service...")
	srv := service.NewMailService(cache, repo, conf.SMTP, conf.Code)

	logrus.Debug("initializing controller...")
	echo := echo.New()
	v1.NewController(echo, srv, infoLog, errLog)

	logrus.Debug("initializing http server...")
	return httpserver.New(echo, conf.HTTPServer).Run(ctx)
}
