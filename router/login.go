package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/Depado/pb-templ-htmx-todo/components"
	"github.com/Depado/pb-templ-htmx-todo/components/shared"
	"github.com/Depado/pb-templ-htmx-todo/htmx"
)

func (ar *AppRouter) Login(c echo.Context, email string, password string) error {
	user, err := ar.App.Dao().FindAuthRecordByEmail("users", email)
	if err != nil {
		return fmt.Errorf("Invalid credentials.")
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

	return components.Render(http.StatusOK, c, components.Login(shared.Context{}, components.LoginFormValue{}, nil, htmx.IsHtmxRequest(c)))
}

func (ar *AppRouter) PostLogin(c echo.Context) error {
	form := components.GetLoginFormValue(c)
	err := form.Validate()

	if err == nil {
		err = ar.Login(c, form.Email, form.Password)
	}

	if err != nil {
		return components.Render(http.StatusOK, c, components.LoginForm(form, err))
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
