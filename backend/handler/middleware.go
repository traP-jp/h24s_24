package handler

import (
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

type userNameContextKey string

const userNameCtxKey = "userName"

func userNameMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Request().Header.Get("X-Forwarded-User")
		if username != "" {
			c.Set(userNameCtxKey, username)
			return next(c)
		}

		local, err := strconv.ParseBool(os.Getenv("LOCAL"))
		if err != nil || !local {
			return echo.NewHTTPError(403, "no username")
		}

		c.Set(userNameCtxKey, "test-user2")
		return next(c)
	}
}
