package router

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/Depado/pb-templ-htmx-todo/components"
	"github.com/Depado/pb-templ-htmx-todo/htmx"
)

func CustomHTTPErrorHandler(c echo.Context, err error) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	c.Response().Writer.WriteHeader(code)
	components.HTTPError(code, "", htmx.IsHtmxRequest(c)).Render(c.Request().Context(), c.Response().Writer)
}
