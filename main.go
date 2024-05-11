package main

import (
	"log"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	// _ "github.com/Depado/pb-templ-htmx-todo/migrations"
	"github.com/Depado/pb-templ-htmx-todo/assets"
	"github.com/Depado/pb-templ-htmx-todo/components"
	"github.com/Depado/pb-templ-htmx-todo/components/httperror"
	"github.com/Depado/pb-templ-htmx-todo/utils"
)

func main() {
	app := pocketbase.New()
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.HTTPErrorHandler = httperror.CustomHTTPErrorHandler
		e.Router.GET("/static/*", assets.AssetsHandler(app.Logger(), isGoRun), middleware.Gzip())

		e.Router.GET("/", func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
			return components.Home(utils.IsHtmxRequest(c)).Render(c.Request().Context(), c.Response().Writer)
		})

		e.Router.GET("/login", func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
			return components.Login(utils.IsHtmxRequest(c)).Render(c.Request().Context(), c.Response().Writer)
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
