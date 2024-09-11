package routes

import (
	"net/http"

	"github.com/dagulv/stock-api/internal/adapter/json"
	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

type userRoutes struct {
	Json    jsoniter.API
	Service service.User
}

func UserRoutes(e *echo.Echo, s service.User, jsonApi jsoniter.API) {
	r := userRoutes{
		Json:    jsonApi,
		Service: s,
	}

	e.GET("/users", r.list)
	e.GET("/users/:id", r.get)
	e.POST("/users", r.create)
}

func (r userRoutes) list(c echo.Context) (err error) {
	domainEncoder := json.CreateDomainEncoder[*domain.User](r.Json, c.Response())
	defer r.Json.ReturnStream(domainEncoder.Stream)

	list, err := r.Service.List(c.Request().Context())

	if err != nil {
		return
	}

	for u := range list {
		domainEncoder.Add(&u)
	}

	return domainEncoder.Flush()
}

func (r userRoutes) get(c echo.Context) (err error) {
	var user domain.User
	var userId xid.ID

	if userId, err = xid.FromString(c.Param("id")); err != nil {
		return
	}

	if err = r.Service.Get(c.Request().Context(), userId, &user); err != nil {
		return
	}

	return c.JSON(http.StatusOK, user)
}

func (r userRoutes) create(c echo.Context) (err error) {
	var user domain.User

	if err = c.Bind(&user); err != nil {
		return
	}

	if err = c.Validate(&user); err != nil {
		return
	}

	if err = r.Service.Create(c.Request().Context(), &user); err != nil {
		return
	}

	return c.JSON(http.StatusOK, user)
}
