package server

import (
	"github.com/essentialkaos/branca/v2"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token branca.Branca

		cookie, err := c.Cookie("session")

		if err != nil {
			return next(c)
		}

		token = branca.DecodeString(cookie.Value)
	}
}
