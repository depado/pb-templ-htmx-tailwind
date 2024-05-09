package assets

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
)

//go:embed all:dist
var assets embed.FS

func AssetsHandler(live bool) echo.HandlerFunc {
	assetHandler := http.FileServer(GetFileSystem(live))
	return echo.WrapHandler(http.StripPrefix("/static/", assetHandler))
}

func GetFileSystem(live bool) http.FileSystem {

	if live {
		log.Print("using live mode")
		return http.FS(os.DirFS("assets/dist"))
	}

	log.Print("using embed mode")
	fsys, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
