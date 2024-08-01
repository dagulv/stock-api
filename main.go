package main

import (
	"context"

	"github.com/dagulv/stock-api/internal/adapter/http"
	"github.com/dagulv/stock-api/internal/adapter/timescale"
	"github.com/dagulv/stock-api/internal/core/service"
	"github.com/dagulv/stock-api/internal/env"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := start(ctx); err != nil {
		panic(err)
	}
}

func start(ctx context.Context) (err error) {
	env, err := env.GetEnv()

	if err != nil {
		return
	}

	db, err := timescale.Connect(ctx, env)

	if err != nil {
		return
	}
	defer db.Close()

	tickerService := service.Ticker{
		Store: timescale.NewTicker(db),
	}

	// wconfig := &webauthn.Config{
	// 	RPDisplayName: "Stock",
	// 	RPID:          "stock.local",
	// 	RPOrigins:     []string{"https://stock.local"},
	// }

	// webAuthn, err := webauthn.New(wconfig)

	// if err != nil {
	// 	return
	// }

	// json := jsoniter.ConfigFastest

	server := http.Server{
		// Json:     json,
		// WebAuthn: webAuthn,
		Tick: tickerService,
	}
	if err = tickerService.Spawn(ctx); err != nil {
		return
	}

	return server.StartServer(ctx)
}
