package server

import (
	"errors"
	"net/http"

	"github.com/danztran/sample-admission-validating-webhook/pkg/httpclient"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *server) setupAPIs(e *echo.Echo) error {
	e.GET("/health", func(c echo.Context) error { return c.String(http.StatusOK, "OK") })
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.POST("/", wrapHandler(handler))

	return nil
}

func wrapHandler(hl func(echo.Context) error) func(echo.Context) error {
	return func(c echo.Context) error {
		err := hl(c)
		if err != nil {
			return catchHandlerError(c, err)
		}
		return nil
	}
}

func catchHandlerError(c echo.Context, err error) error {
	var errNotFound *httpclient.ErrNotFound

	switch true {
	case errors.As(err, &errNotFound):
		err = c.String(http.StatusNotFound, err.Error())
	default:
		c.Error(err)
	}

	return err
}
