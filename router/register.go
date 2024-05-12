package router

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/models"
)

func (ar *AppRouter) Register(c echo.Context, email string, password string, passwordRepeat string) error {
	user, _ := ar.App.Dao().FindAuthRecordByEmail("users", email)
	if user != nil {
		return fmt.Errorf("username already taken")
	}

	if password != passwordRepeat {
		return fmt.Errorf("passwords don't match")
	}

	collection, err := ar.App.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}

	newUser := models.NewRecord(collection)
	newUser.SetPassword(password)
	newUser.SetEmail(email)

	if err = ar.App.Dao().SaveRecord(newUser); err != nil {
		return err
	}

	return ar.setAuthToken(c, newUser)
}
