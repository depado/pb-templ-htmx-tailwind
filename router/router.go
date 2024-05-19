package router

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase/core"

	"github.com/Depado/pb-templ-htmx-tailwind/assets"
	"github.com/Depado/pb-templ-htmx-tailwind/htmx"
)

type AppRouter struct {
	App    core.App
	Router *echo.Echo
}

func NewAppRouter(e *core.ServeEvent) *AppRouter {
	return &AppRouter{
		App:    e.App,
		Router: e.Router,
	}
}

func (ar *AppRouter) SetupRoutes(live bool) error {
	ar.Router.Use(middleware.Logger())
	ar.Router.HTTPErrorHandler = htmx.WrapDefaultErrorHandler(ar.Router.HTTPErrorHandler)
	ar.Router.GET("/static/*", assets.AssetsHandler(ar.App.Logger(), live), middleware.Gzip())

	ar.Router.Use(ar.LoadAuthContextFromCookie())
	ar.Router.GET("/", ar.GetHome)
	ar.Router.GET("/login", ar.GetLogin)
	ar.Router.POST("/login", ar.PostLogin)
	ar.Router.POST("/register", ar.PostRegister)
	ar.Router.POST("/logout", ar.PostLogout)
	ar.Router.GET("/error", ar.GetError)

	ar.Router.PATCH("/task/:id", ar.ToggleTask)
	ar.Router.POST("/list/:id/task", ar.CreateTask)

	// unimplemented
	// ar.Router.POST("/list", ar.ListCreate)
	// ar.Router.PATCH("/list/:id/archive", ar.ListArchive)
	// ar.Router.DELETE("/list/:id", ar.ListDelete)

	return nil
}
