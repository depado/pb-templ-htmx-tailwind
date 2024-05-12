package components

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

// Render a templ component with the appropriate headers.
func Render(code int, ctx echo.Context, c templ.Component) error {
	ctx.Response().Writer.WriteHeader(code)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return c.Render(ctx.Request().Context(), ctx.Response().Writer)
}
