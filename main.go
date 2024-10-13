package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/dagulv/stock-api/internal/adapter/db"
	"github.com/dagulv/stock-api/internal/adapter/http"
	"github.com/dagulv/stock-api/internal/adapter/mailer"
	"github.com/dagulv/stock-api/internal/adapter/server"
	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/port"
	"github.com/dagulv/stock-api/internal/core/service"
	"github.com/dagulv/stock-api/internal/env"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/jackc/pgx/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailersend/mailersend-go"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
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
	parser.AddRule(paseto.NotExpired())

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

	stockService := service.Stock{
		Store: db.NewStock(dbPool),
	}

	wconfig := &webauthn.Config{
		RPDisplayName: "Stock",
		RPID:          "stock.local",
		RPOrigins:     []string{"https://stock.local"},
	}

	webAuthn, err := webauthn.New(wconfig)

	if err != nil {
		return
	}

	json := jsoniter.ConfigFastest

	if err != nil {
		return
	}

	server := http.Server{
		Json:     json,
		WebAuthn: webAuthn,
		Auth:     authService,
		Tick:     tickerService,
		User:     userService,
		Stock:    stockService,
	}

	if err = createUser(ctx, authService.Store, userService.Store, env); err != nil {
		return
	}

	if err = tickerService.Spawn(ctx); err != nil {
		return
	}

	return server.StartServer(ctx)
}

func createUser(ctx context.Context, authStore port.Auth, userStore port.User, env env.Env) (err error) {
	credParts := strings.Split(env.AdminUserCred, ":")

	if len(credParts) < 2 {
		return errors.New("need admin email and password separated by :")
	}

	var existingCreds domain.Credentials

	err = authStore.GetCredentialsByEmail(ctx, credParts[0], &existingCreds)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return
		}

	}

	if !existingCreds.UserId.IsNil() {
		return nil
	}

	now := time.Now()
	user := domain.User{
		Id:          xid.New(),
		Email:       credParts[0],
		TimeCreated: now,
		TimeUpdated: now,
	}

	if err = userStore.Create(ctx, &user); err != nil {
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credParts[1]), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	creds := domain.Credentials{
		UserId:   user.Id,
		Password: string(hashedPassword),
	}

	if err = authStore.InsertCredentials(ctx, creds); err != nil {
		return
	}

	return nil
}
