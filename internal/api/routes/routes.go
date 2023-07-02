package routes

import "github.com/labstack/echo/v4"

type Routes interface {
	CurrentRoutes(*echo.Group)
}
