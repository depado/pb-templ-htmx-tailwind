package auth

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v5"
)

const (
	EmailField          = "email"
	UsernameField       = "username"
	PasswordRepeatField = "password_repeat"
)

type RegisterFormValue struct {
	Email          string
	Username       string
	Password       string
	PasswordRepeat string
}

type RegisterFormErrors struct {
	Email          error
	Username       error
	Password       error
	PasswordRepeat error
}

func (rfv RegisterFormValue) Validate() (RegisterFormErrors, error) {
	rfe := RegisterFormErrors{
		Email:          validation.Validate(rfv.Email, validation.Required.Error("Required"), is.Email.Error("Invalid email address")),
		Username:       validation.Validate(rfv.Username, validation.Required.Error("Required"), validation.Length(5, 50).Error("Length must be between 5 and 50 chars")),
		Password:       validation.Validate(rfv.Password, validation.Required.Error("Required"), validation.Length(10, 50).Error("Length must be between 10 and 50 chars")),
		PasswordRepeat: validation.Validate(rfv.PasswordRepeat, validation.Required.Error("Required")),
	}

	if rfv.Password != "" && rfv.PasswordRepeat != "" && rfv.Password != rfv.PasswordRepeat {
		rfe.PasswordRepeat = fmt.Errorf("Passwords don't match")
	}

	if rfe.Email != nil || rfe.Username != nil || rfe.Password != nil || rfe.PasswordRepeat != nil {
		return rfe, fmt.Errorf("Validation error")
	}

	return rfe, nil
}

func GetRegisterFormValue(c echo.Context) RegisterFormValue {
	return RegisterFormValue{
		Email:          c.FormValue(EmailField),
		Username:       c.FormValue(UsernameField),
		Password:       c.FormValue(PasswordField),
		PasswordRepeat: c.FormValue(PasswordRepeatField),
	}
}
