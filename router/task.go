package router

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/Depado/pb-templ-htmx-tailwind/components"
	"github.com/Depado/pb-templ-htmx-tailwind/htmx"
	"github.com/Depado/pb-templ-htmx-tailwind/models"
)

func (ar *AppRouter) ToggleTask(c echo.Context) error {
	rec := c.Get(apis.ContextAuthRecordKey)
	if rec == nil {
		return htmx.Redirect(c, "/")
	}

	id := c.PathParam("id")
	task, err := models.GetTaskById(ar.App.Dao(), id)
	if err != nil {
		ar.App.Logger().Debug("get task by id error", "error", err, "task", id)
		return htmx.Error(c, "Unable to find this task")
	}

	task.Done = !task.Done

	if err := task.Save(ar.App.Dao()); err != nil {
		ar.App.Logger().Error("saving task", "error", err, "task", id)
		return htmx.Error(c, "Unable to save this task")
	}

	list, err := models.GetListById(ar.App.Dao(), task.ListId, true)
	if err != nil {
		ar.App.Logger().Error("get list by id error", "error", err, "list", task.ListId)
		return htmx.Error(c, "Unable to find this list")
	}

	return components.Render(c, http.StatusOK, components.List(list))
}

func (ar *AppRouter) CreateTask(c echo.Context) error {
	rec := c.Get(apis.ContextAuthRecordKey)
	if rec == nil {
		return htmx.Redirect(c, "/")
	}

	t := &models.Task{
		Title:  c.FormValue("title"),
		ListId: c.PathParam("id"),
	}

	if err := t.Save(ar.App.Dao()); err != nil {
		ar.App.Logger().Error("save new task", "error", err, "list", t.ListId)
		return htmx.Error(c, "Unable to save new task")
	}

	list, err := models.GetListById(ar.App.Dao(), t.ListId, true)
	if err != nil {
		ar.App.Logger().Error("get list by id error", "error", err, "list", t.ListId)
		return htmx.Error(c, "Unable to find this list")
	}

	return components.Render(c, http.StatusOK, components.List(list))
}
