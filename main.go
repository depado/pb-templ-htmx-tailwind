package main

import (
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "github.com/depado/pb-templ-htmx-tailwind/migrations"
	"github.com/depado/pb-templ-htmx-tailwind/router"
)

func main() {
	app := pocketbase.New()
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return router.NewAppRouter(e).SetupRoutes(isGoRun)
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
