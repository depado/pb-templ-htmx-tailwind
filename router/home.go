package router

import (
	"net/http"

	"github.com/Depado/pb-templ-htmx-tailwind/components"
	"github.com/Depado/pb-templ-htmx-tailwind/components/shared"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
)

func GetHome(c echo.Context) error {
	var user *models.Record
	rec := c.Get(apis.ContextAuthRecordKey)
	if rec != nil {
		user = c.Get(apis.ContextAuthRecordKey).(*models.Record)
	}
	return components.Render(http.StatusOK, c, components.Home(shared.Context{User: user}))
}
