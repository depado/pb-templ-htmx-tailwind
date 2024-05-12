package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

func (ar *AppRouter) setAuthToken(c echo.Context, user *models.Record) error {
	s, tokenErr := tokens.NewRecordAuthToken(ar.App, user)
	if tokenErr != nil {
		return fmt.Errorf("login failed")
	}

	c.SetCookie(&http.Cookie{
		Name:     AuthCookieName,
		Value:    s,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	return nil
}
