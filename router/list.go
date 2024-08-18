package router

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/depado/pb-templ-htmx-tailwind/components"
	"github.com/depado/pb-templ-htmx-tailwind/htmx"
	"github.com/depado/pb-templ-htmx-tailwind/models"
)

func (ar *AppRouter) ToggleArchive(c echo.Context) error {
	rec := c.Get(apis.ContextAuthRecordKey)
	if rec == nil {
		return htmx.Redirect(c, "/")
	}
	id := c.PathParam("id")

	l, err := models.GetListById(ar.App.Dao(), id, false)
	if err != nil {
		ar.App.Logger().Error("toggle archive: get list by id", "error", err, "id", id)
		return htmx.Error(c, "Unable to find this list")
	}

	l.Archived = !l.Archived

	if err := l.Save(ar.App.Dao()); err != nil {
		ar.App.Logger().Error("toggle archive: save list", "error", err, "id", id)
		return htmx.Error(c, "Unable to save list")
	}

	msg := "List archived"
	if !l.Archived {
		msg = "List unarchived"
	}

	return components.Render(c, http.StatusOK, components.ListWithToast(l, msg))
}

func (ar *AppRouter) ListDelete(c echo.Context) error {
	rec := c.Get(apis.ContextAuthRecordKey)
	if rec == nil {
		return htmx.Redirect(c, "/")
	}
	id := c.PathParam("id")

	l, err := models.GetListById(ar.App.Dao(), id, false)
	if err != nil {
		ar.App.Logger().Error("delete list: get list by id", "error", err, "id", id)
		return htmx.Error(c, "Unable to find this list")
	}

	if err := ar.App.Dao().Delete(l); err != nil {
		ar.App.Logger().Error("delete list: delete:", "error", err, "id", id)
		return htmx.Error(c, "Unable to delete this list")
	}

	return htmx.Redirect(c, "/")
}
