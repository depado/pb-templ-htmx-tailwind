package components

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
)

type LoginFormValue struct {
	Email    string
	Password string
}

func (lfv LoginFormValue) Validate() error {
	return validation.ValidateStruct(&lfv,
		validation.Field(&lfv.Email, validation.Required, validation.Length(3, 50)),
		validation.Field(&lfv.Password, validation.Required),
	)
}

func GetLoginFormValue(c echo.Context) LoginFormValue {
	return LoginFormValue{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}
}
