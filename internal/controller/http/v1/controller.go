package v1

import (
	"github.com/cutlery47/email-service/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/cutlery47/email-service/docs"
)

func NewController(e *echo.Echo, srv service.Service, infoLog, errLog *logrus.Logger) {
	e.Use(middleware.Recover())

	// healthcheck
	e.GET("/ping", func(c echo.Context) error { return c.NoContent(200) })
	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("/api/v1", requestLoggerMiddleware(infoLog))
	{
		newMailRoutes(v1, srv, newErrMapper(errLog))
	}

}
