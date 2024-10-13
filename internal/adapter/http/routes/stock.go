package routes

import (
	"net/http"

	"github.com/dagulv/stock-api/internal/adapter/json"
	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type stockRoutes struct {
	Json    jsoniter.API
	Service service.Stock
}

func StockRoutes(e *echo.Echo, s service.Stock, jsonApi jsoniter.API) {
	r := stockRoutes{
		Json:    jsonApi,
		Service: s,
	}

	e.GET("/stocks", r.list)
	e.GET("/stocks/:symbol", r.get)
}

func (r stockRoutes) list(c echo.Context) (err error) {
	domainEncoder := json.CreateDomainEncoder[*domain.Stock](r.Json, c.Response())
	defer r.Json.ReturnStream(domainEncoder.Stream)

	list, err := r.Service.List(c.Request().Context())

	if err != nil {
		return
	}

	for row := range list {
		domainEncoder.Add(&row)
	}

	return domainEncoder.Flush()
}

func (r stockRoutes) get(c echo.Context) (err error) {
	var stock domain.Stock

	if err = r.Service.Get(c.Request().Context(), c.Param("symbol"), &stock); err != nil {
		return
	}

	return c.JSON(http.StatusOK, stock)
}
