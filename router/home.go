package router

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	pbmodels "github.com/pocketbase/pocketbase/models"

	"github.com/depado/pb-templ-htmx-tailwind/components"
	"github.com/depado/pb-templ-htmx-tailwind/components/shared"
	"github.com/depado/pb-templ-htmx-tailwind/htmx"
	"github.com/depado/pb-templ-htmx-tailwind/models"
)

func (ar *AppRouter) GetHome(c echo.Context) error {
	rec := c.Get(apis.ContextAuthRecordKey)
	if rec == nil {
		return components.Render(c, http.StatusOK, components.Home(shared.Context{}, false))
	}

	user := c.Get(apis.ContextAuthRecordKey).(*pbmodels.Record)
	lists, err := models.FindUserLists(ar.App.Dao(), user.Id)
	if err != nil {
		ar.App.Logger().Error("unable to get lists for user", "error", err, "id", user.Id)
		return htmx.Error(c, "Unable to get lists")
	}

	return components.Render(c, http.StatusOK, components.Home(shared.Context{User: user, Lists: lists}, false))
}
