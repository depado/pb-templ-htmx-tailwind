package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/Depado/pb-templ-htmx-tailwind/components"
	"github.com/Depado/pb-templ-htmx-tailwind/components/auth"
	"github.com/Depado/pb-templ-htmx-tailwind/components/shared"
	"github.com/Depado/pb-templ-htmx-tailwind/htmx"
)

func (ar *AppRouter) Login(c echo.Context, identifier string, password string) error {
	user, err := ar.App.Dao().FindAuthRecordByEmail("users", identifier)
	if err != nil {
		user, err = ar.App.Dao().FindAuthRecordByUsername("users", identifier)
		if err != nil {
			return fmt.Errorf("Invalid credentials.")
		}
	}

	valid := user.ValidatePassword(password)
	if !valid {
		return fmt.Errorf("Invalid credentials.")
	}

	return ar.setAuthToken(c, user)
}

func (ar *AppRouter) GetLogin(c echo.Context) error {
	if c.Get(apis.ContextAuthRecordKey) != nil {
		return c.Redirect(302, "/")
	}

	return components.Render(http.StatusOK, c, auth.LoginPage(shared.Context{}, auth.LoginPageForms{}))
}

func (ar *AppRouter) PostLogin(c echo.Context) error {
	form := auth.GetLoginFormValue(c)
	lpe, err := form.Validate()

	if err == nil {
		err = ar.Login(c, form.Identifier, form.Password)
	}

	if err != nil {
		return components.Render(http.StatusOK, c, auth.LoginForm(form, lpe, err))
	}

	return htmx.Redirect(c, "/")
}

func (ar *AppRouter) PostLogout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     AuthCookieName,
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
	})

	return htmx.Redirect(c, "/")
}
