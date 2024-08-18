package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/depado/pb-templ-htmx-tailwind/components"
	"github.com/depado/pb-templ-htmx-tailwind/components/auth"
	"github.com/depado/pb-templ-htmx-tailwind/components/shared"
	"github.com/depado/pb-templ-htmx-tailwind/htmx"
)

// Login handles the login logic.
func (ar *AppRouter) Login(c echo.Context, identifier string, password string) error {

	// Try to find the user by email first, then by username
	user, err := ar.App.Dao().FindAuthRecordByEmail("users", identifier)
	if err != nil {
		user, err = ar.App.Dao().FindAuthRecordByUsername("users", identifier)
		if err != nil {
			ar.App.Logger().Error("unknown user", "identifier", identifier)
			return fmt.Errorf("Invalid credentials.")
		}
	}

	// Actually verify password
	valid := user.ValidatePassword(password)
	if !valid {
		ar.App.Logger().Error("wrong password", "identifier", identifier, "id", user.Id)
		return fmt.Errorf("Invalid credentials.")
	}

	ar.App.Logger().Debug("user signed-in", "identifier", identifier, "id", user.Id)
	return ar.setAuthToken(c, user)
}

// GetLogin returns the login page.
func (ar *AppRouter) GetLogin(c echo.Context) error {
	if c.Get(apis.ContextAuthRecordKey) != nil {
		return c.Redirect(302, "/")
	}

	return components.Render(c, http.StatusOK, auth.LoginPage(shared.Context{}, auth.LoginPageForms{}))
}

// PostLogin handles the form validation and authentication.
func (ar *AppRouter) PostLogin(c echo.Context) error {
	form := auth.GetLoginFormValue(c)
	lpe, err := form.Validate()

	if err == nil {
		err = ar.Login(c, form.Identifier, form.Password)
	}

	if err != nil {
		return components.Render(c, http.StatusOK, auth.LoginForm(form, lpe, err))
	}

	return htmx.Redirect(c, "/")
}

// PostLogout logs the user out by clearing the authentication cookie.
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
