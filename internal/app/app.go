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

//	@title			Email Service
//	@version		0.0.1
//	@description	This is an email authentication service

//	@contact.name	Ivanchenko Arkhip
//	@contact.email	kitchen_cutlery@mail.ru

//	@BasePath	/

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

	var cache service.Cache

	logrus.Debug("initializing cache...")
	if conf.CacheType == "redis" {
		cache, err = service.NewRedisCache(ctx, conf.Cache)
		if err != nil {
			return fmt.Errorf("error when setting up redis: %v", err)
		}
	} else if conf.CacheType == "map" {
		cache = service.NewMapCache(conf.Cache, infoLog)
	} else {
		return fmt.Errorf("valid cache types: redis, map")
	}

	logrus.Debug("initializing repository...")
	repo, err := repo.NewMailRepository(ctx, conf.Postgres, infoLog)
	if err != nil {
		return fmt.Errorf("error when connecting to the database: %v", err)
	}

	logrus.Debug("initializing service...")
	srv := service.NewMailService(cache, repo, conf.SMTP, conf.Code)

	logrus.Debug("initializing controller...")
	echo := echo.New()
	v1.NewController(echo, srv, infoLog, errLog)

	logrus.Debug("initializing http server...")
	return httpserver.New(echo, conf.HTTPServer).Run(ctx)
}
