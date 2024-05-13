package router

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

func (ar *AppRouter) GetError(c echo.Context) error {
	return fmt.Errorf("oh no")
}
