package router

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

const AuthCookieName = "Auth"

func (ar *AppRouter) LoadAuthContextFromCookie() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenCookie, err := c.Request().Cookie(AuthCookieName)
			if err != nil || tokenCookie.Value == "" {
				return next(c) // no token cookie
			}

			token := tokenCookie.Value

			claims, _ := security.ParseUnverifiedJWT(token)
			tokenType := cast.ToString(claims["type"])

			switch tokenType {
			case tokens.TypeAdmin:
				admin, err := ar.App.Dao().FindAdminByToken(
					token,
					ar.App.Settings().AdminAuthToken.Secret,
				)
				if err == nil && admin != nil {
					c.Set(apis.ContextAdminKey, admin)
				}

			case tokens.TypeAuthRecord:
				record, err := ar.App.Dao().FindAuthRecordByToken(
					token,
					ar.App.Settings().RecordAuthToken.Secret,
				)
				if err == nil && record != nil {
					c.Set(apis.ContextAuthRecordKey, record)
				}
			}

			return next(c)
		}
	}
}
