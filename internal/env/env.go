package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseUrl string
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
	return
}

func LoadEnv(path string) (err error) {
	return godotenv.Load(path)
}
