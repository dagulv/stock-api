package routes

import (
	"log"
	"net/http"

	"github.com/dagulv/stock-api/internal/models"
	"github.com/dagulv/stock-api/internal/services"
	"github.com/dagulv/stock-api/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

type UserRoutes struct {
	Service *services.UserService
}

func (r *UserRoutes) CurrentRoutes(e *echo.Group) {
	e.GET("/users", r.list)
	e.GET("/users/:id", r.get)
	e.POST("/users", r.create)
}

func (r *UserRoutes) list(c echo.Context) (err error) {
	stream := utils.InitEncoder[*models.User](c.Response())
	defer stream.EndStream()

	if err = r.Service.List(c.Request().Context(), stream.Iterator); err != nil {
		return
	}

	return nil
}

func (r *UserRoutes) get(c echo.Context) (err error) {
	id := c.Param("id")
	var user models.User

	if id == "me" {
		// id = string(xid.New())
	}

	userId, err := xid.FromString(id)

	if err != nil {
		return
	}

	if err = r.Service.Get(c.Request().Context(), userId, &user); err != nil {
		log.Println(err)
		return
	}

	return c.JSON(http.StatusOK, user)
}

func (r *UserRoutes) create(c echo.Context) (err error) {
	user := models.User{}

	if err = c.Bind(&user); err != nil {
		return
	}

	r.Service.Create(c.Request().Context(), &user)

	return c.JSON(http.StatusOK, user)
}
