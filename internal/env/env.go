package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseUrl      string
	AuthSecretKey    string
	AppUrl           string
	MailerSendApiKey string
	RecipientName    string
	RecipientEmail   string
	AdminUserCred    string
}

func GetEnv(paths ...string) (e Env, err error) {
	var path string

	if len(paths) > 0 {
		path = paths[0]
	}

	if path == "" {
		err = godotenv.Load()
	} else {
		err = godotenv.Load(path)
	}

	if err != nil {
		return
	}

	e.DatabaseUrl = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))

	e.AuthSecretKey = os.Getenv("AUTH_SECRET_KEY")

	e.AppUrl = os.Getenv("APP_URL")

	e.MailerSendApiKey = os.Getenv("MAILERSEND_API_KEY")

	e.RecipientName = os.Getenv("RECIPIENT_NAME")
	e.RecipientEmail = os.Getenv("RECIPIENT_EMAIL")

	e.AdminUserCred = os.Getenv("ADMIN_USER")

	return
}

func LoadEnv(path string) (err error) {
	return godotenv.Load(path)
}
