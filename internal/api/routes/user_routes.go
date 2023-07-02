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
	e.DELETE("/users", r.delete)
	e.PUT("/users/:id", r.put)
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
	var user models.User

	if err = c.Bind(&user); err != nil {
		return
	}

	r.Service.Create(c.Request().Context(), &user)

	return c.JSON(http.StatusOK, user)
}

func (r *UserRoutes) put(c echo.Context) (err error) {
	id, err := getId(c)

	if err != nil {
		return
	}

	var user models.User

	if err = c.Bind(&user); err != nil {
		return
	}

	user.Id = id

	if err = r.Service.Put(c.Request().Context(), &user); err != nil {
		return
	}

	return c.JSON(http.StatusOK, &user)
}

func (r *UserRoutes) delete(c echo.Context) (err error) {
	userId, err := getId(c)

	if err != nil {
		return
	}

	if err = r.Service.Delete(c.Request().Context(), userId); err != nil {
		return
	}

	return c.NoContent(http.StatusOK)
}

func getId(c echo.Context) (userId xid.ID, err error) {
	id := c.Param("id")

	userId, err = xid.FromString(id)

	return
}
