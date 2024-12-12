package v1

import (
	"github.com/cutlery47/email-service/internal/service"
	"github.com/labstack/echo/v4"
)

type mailRoutes struct {
	srv service.Service
	e   *errMapper
}

func newMailRoutes(g *echo.Group, srv service.Service, e *errMapper) {
	r := &mailRoutes{
		srv: srv,
		e:   e,
	}

	g.GET("/register", r.registerUser)
	g.GET("/confirm", r.confirmRegister)
}

func (mr *mailRoutes) registerUser(c echo.Context) error {
	return nil
}

func (mr *mailRoutes) confirmRegister(c echo.Context) error {
	return nil
}
