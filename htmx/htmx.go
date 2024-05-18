package htmx

import (
	"net/http"
	"strings"

	"github.com/Depado/pb-templ-htmx-tailwind/components"
	"github.com/Depado/pb-templ-htmx-tailwind/components/shared"
	"github.com/labstack/echo/v5"
)

// IsHtmxRequest checks if the received request has the HX-Request header that
// indicates a request was performed by HTMX.
func IsHtmxRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}

// Redirect handles redirection properly for HTMX.
func Redirect(c echo.Context, path string) error {
	if IsHtmxRequest(c) {
		c.Response().Header().Set("HX-Location", path)
		return c.NoContent(204)
	}

	return c.Redirect(302, path)
}

// WrapDefaultErrorHandler wraps the provided error handler to properly serve
// HTML for HTMX and falls back to the provided error handler for other requests
// including those to Pocketbase.
func WrapDefaultErrorHandler(defaultErrorHandler echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(c echo.Context, err error) {
		if IsHtmxRequest(c) || (!strings.HasPrefix(c.Path(), "/_/") && !strings.HasPrefix(c.Path(), "/api/")) {
			code := http.StatusInternalServerError
			if he, ok := err.(*echo.HTTPError); ok {
				code = he.Code
			}

			components.Render(code, c, components.HTTPError(shared.Context{}, code, "", IsHtmxRequest(c))) //nolint:errcheck
		} else {
			defaultErrorHandler(c, err)
		}
	}
}
