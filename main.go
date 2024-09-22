package main

import (
	"context"

	"aidanwoods.dev/go-paseto"
	"github.com/dagulv/stock-api/internal/adapter/db"
	"github.com/dagulv/stock-api/internal/adapter/http"
	"github.com/dagulv/stock-api/internal/adapter/mailer"
	"github.com/dagulv/stock-api/internal/adapter/server"
	"github.com/dagulv/stock-api/internal/core/service"
	"github.com/dagulv/stock-api/internal/env"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailersend/mailersend-go"
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

	secretKey := server.GenerateSecretKeyFromEnv(env)

	dbPool, err := db.Connect(ctx, env)

	if err != nil {
		return
	}

	defer dbPool.Close()

	parser := paseto.NewParser()
	parser.AddRule(paseto.Subject("id"), paseto.NotExpired())

	ms := mailersend.NewMailersend(env.MailerSendApiKey)
	mailer := mailer.New(ctx, mailer.Mailersend{
		From: mailersend.Recipient{
			Name:  env.RecipientName,
			Email: env.RecipientEmail,
		},
		Ms: ms,
	}, 10)

	authService := service.Auth{
		SecretAuthKey: secretKey,
		Parser:        parser,
		Mailer:        mailer,
		Env:           env,
		Store:         db.NewAuth(dbPool),
	}

	tickerService := service.Ticker{
		Store: db.NewTicker(dbPool),
	}

	userService := service.User{
		Store: db.NewUser(dbPool),
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

	json := jsoniter.ConfigFastest

	if err != nil {
		return
	}

	server := http.Server{
		Json: json,
		// WebAuthn: webAuthn,
		Auth: authService,
		Tick: tickerService,
		User: userService,
	}
	if err = tickerService.Spawn(ctx); err != nil {
		return
	}

	return server.StartServer(ctx)
}
