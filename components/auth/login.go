package auth

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
)

const (
	IdentifierField = "email"
	PasswordField   = "password"
)

type LoginFormValue struct {
	Identifier string
	Password   string
}

type LoginFormErrors struct {
	Identifier error
	Password   error
}

func (lfv LoginFormValue) Validate() (LoginFormErrors, error) {
	lfe := LoginFormErrors{
		Identifier: validation.Validate(
			lfv.Identifier,
			validation.Required.Error("Required"),
			validation.Length(5, 50).Error("Length must be between 5 and 50 chars"),
		),
		Password: validation.Validate(lfv.Password, validation.Required.Error("Required")),
	}
	if lfe.Identifier != nil || lfe.Password != nil {
		return lfe, fmt.Errorf("validation error")
	}

	return lfe, nil
}

func GetLoginFormValue(c echo.Context) LoginFormValue {
	return LoginFormValue{
		Identifier: c.FormValue(IdentifierField),
		Password:   c.FormValue(PasswordField),
	}
}
