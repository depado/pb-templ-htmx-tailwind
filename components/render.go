package components

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

// Render renders a templ component with the appropriate headers.
func Render(ctx echo.Context, code int, c templ.Component) error {
	ctx.Response().Writer.WriteHeader(code)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return c.Render(ctx.Request().Context(), ctx.Response().Writer)
}
