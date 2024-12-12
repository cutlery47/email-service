package v1

import (
	"github.com/cutlery47/email-service/internal/repo"
	"github.com/cutlery47/email-service/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var errMap = map[error]*echo.HTTPError{
	service.ErrCacheNotFound: echo.ErrNotFound,
	service.ErrWrongCode:     echo.ErrBadRequest,
	repo.ErrAlreadyExists:    echo.ErrBadRequest,
}

type errMapper struct {
	errLog *logrus.Logger
}

func newErrMapper(errLog *logrus.Logger) *errMapper {
	return &errMapper{
		errLog: errLog,
	}
}

func (e errMapper) Map(err error) *echo.HTTPError {
	if httpErr, ok := errMap[err]; ok {
		httpErr.Message = err.Error()
		return httpErr
	}

	e.errLog.Error(err.Error())
	return echo.ErrInternalServerError
}
