package htmx

import "github.com/labstack/echo/v5"

func IsHtmxRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}

func Redirect(c echo.Context, path string) error {
	if IsHtmxRequest(c) {
		c.Response().Header().Set("HX-Location", path)
		return c.NoContent(204)
	}

	return c.Redirect(302, path)
}
