package router

import (
	"net/http"

	"github.com/Depado/pb-templ-htmx-todo/components"
	"github.com/labstack/echo/v5"
)

func GetHome(c echo.Context) error {
	return components.Render(http.StatusOK, c, components.Home())
}
