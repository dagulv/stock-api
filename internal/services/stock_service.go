package services

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dagulv/stock-api/internal/stores"
	"github.com/go-playground/validator/v10"
)

const (
	// YFINURL            = "https://query1.finance.yahoo.com"
	YFINURL            = "http://pokeapi.co/api/v2/pokedex/kanto/"
	defaultHTTPTimeout = 80 * time.Second
)

type StockService struct {
	Store      stores.StockStore
	Validate   *validator.Validate
	HTTPClient *http.Client
}

func (s StockService) Get(ctx context.Context) (err error) {
	req, err := http.NewRequest(http.MethodGet, YFINURL, nil)

	if err != nil {
		return
	}

	res, err := s.HTTPClient.Do(req)

	if err != nil {
		return
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	log.Println(respBody)

	return
}
