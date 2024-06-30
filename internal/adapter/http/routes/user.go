package routes

// import (
// 	"net/http"

// 	"github.com/dagulv/train-api/internal/adapter/json"
// 	"github.com/dagulv/train-api/internal/domain/user"
// 	jsoniter "github.com/json-iterator/go"
// 	"github.com/labstack/echo/v4"
// 	"github.com/rs/xid"
// )

// type userRoute struct {
// 	Json    jsoniter.API
// 	Service user.Service
// }

// func Routes(e *echo.Echo, s user.Service, jsonApi jsoniter.API) {
// 	r := userRoute{
// 		Json:    jsonApi,
// 		Service: s,
// 	}

// 	e.GET("/users", r.list)
// 	e.GET("/users/:id", r.get)
// }

// func (r userRoute) list(c echo.Context) (err error) {
// 	domainEncoder := json.CreateDomainEncoder[*user.User](r.Json, c.Response())
// 	defer r.Json.ReturnStream(domainEncoder.Stream)

// 	if err = r.Service.List(c.Request().Context(), domainEncoder.AddLine); err != nil {
// 		return
// 	}

// 	return domainEncoder.Flush()

// 	// if err = domainEncoder.Flush(); err != nil {
// 	// 	return
// 	// }

// 	// return c.JSONBlob(http.StatusOK, domainEncoder.Stream.Buffer())
// }

// func (r userRoute) get(c echo.Context) (err error) {
// 	var user user.User
// 	var userId xid.ID

// 	if userId, err = xid.FromString(c.Param("id")); err != nil {
// 		return
// 	}

// 	if err = r.Service.Get(c.Request().Context(), userId, &user); err != nil {
// 		return
// 	}

// 	return c.JSON(http.StatusOK, user)
// }
