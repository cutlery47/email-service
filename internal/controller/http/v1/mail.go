package v1

import (
	"encoding/json"
	"io"

	"github.com/cutlery47/email-service/internal/models"
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

	g.POST("/register", r.registerUser)
	g.POST("/confirm", r.confirmRegister)
}

// @Summary		Register User
// @Tags		Email
// @Param		json		body		models.UserData  true "json"
// @Success	200	{object}	string
// @Failure	400	{object}	echo.HTTPError
// @Failure	404	{object}	echo.HTTPError
// @Failure	500	{object}	echo.HTTPError
// @Router		/api/v1/register [post]
func (mr *mailRoutes) registerUser(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		mr.e.errLog.Error(err)
		return echo.NewHTTPError(500, "internal server error")
	}

	user := models.UserData{}

	if err := json.Unmarshal(body, &user); err != nil {
		mr.e.errLog.Error(err)
		return echo.NewHTTPError(400, "invalid json")
	}

	if user.Mail == "" {
		return echo.NewHTTPError(400, "mail must not be empty")
	}

	if err := mr.srv.Register(ctx, user); err != nil {
		return mr.e.Map(err)
	}

	return c.JSON(200, "Confirmation code has been sent to your email!")

}

// @Summary		Confirm Email
// @Tags		Email
// @Param		json		body		models.ConfirmationData  true "json"
// @Success	200	{object}	string
// @Failure	400	{object}	echo.HTTPError
// @Failure	404	{object}	echo.HTTPError
// @Failure	500	{object}	echo.HTTPError
// @Router		/api/v1/confirm [post]
func (mr *mailRoutes) confirmRegister(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		mr.e.errLog.Error(err)
		return echo.NewHTTPError(500, "internal server error")
	}

	confirm := models.ConfirmationData{}

	if err := json.Unmarshal(body, &confirm); err != nil {
		mr.e.errLog.Error(err)
		return echo.NewHTTPError(400, "invalid json")
	}

	if confirm.Code == "" {
		return echo.NewHTTPError(400, "code must not be empty")
	}

	if confirm.Mail == "" {
		return echo.NewHTTPError(400, "mail must not be empty")
	}

	if err := mr.srv.Confirm(ctx, confirm); err != nil {
		return mr.e.Map(err)
	}

	return c.JSON(200, "Your email has been successfully confirmed")
}
