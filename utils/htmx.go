package utils

import "github.com/labstack/echo/v5"

func IsHtmxRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}
