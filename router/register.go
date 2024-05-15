package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/models"

	"github.com/Depado/pb-templ-htmx-tailwind/components"
	"github.com/Depado/pb-templ-htmx-tailwind/components/auth"
	"github.com/Depado/pb-templ-htmx-tailwind/htmx"
)

func (ar *AppRouter) Register(c echo.Context, email, username, password, passwordRepeat string) error {
	user, _ := ar.App.Dao().FindAuthRecordByEmail("users", email)
	if user != nil {
		return fmt.Errorf("Email or username already taken")
	}

	user, _ = ar.App.Dao().FindAuthRecordByUsername("users", username)
	if user != nil {
		return fmt.Errorf("Email or username already taken")
	}

	if password != passwordRepeat {
		return fmt.Errorf("Passwords don't match")
	}

	collection, err := ar.App.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}

	newUser := models.NewRecord(collection)
	if err := newUser.SetPassword(password); err != nil {
		ar.App.Logger().Error("setting password failed", "error", err)
		return fmt.Errorf("Internal error")
	}
	if err := newUser.SetEmail(email); err != nil {
		ar.App.Logger().Error("setting email failed", "error", err)
		return fmt.Errorf("Internal error")
	}
	if err := newUser.SetUsername(username); err != nil {
		ar.App.Logger().Error("setting username failed", "error", err)
		return fmt.Errorf("Internal error")
	}

	if err = ar.App.Dao().SaveRecord(newUser); err != nil {
		return err
	}

	return ar.setAuthToken(c, newUser)
}

func (ar *AppRouter) PostRegister(c echo.Context) error {
	form := auth.GetRegisterFormValue(c)
	rfe, err := form.Validate()

	if err == nil {
		err = ar.Register(c, form.Email, form.Username, form.Password, form.PasswordRepeat)
	}

	if err != nil {
		return components.Render(http.StatusOK, c, auth.RegisterForm(form, rfe, err))
	}

	return htmx.Redirect(c, "/")
}
