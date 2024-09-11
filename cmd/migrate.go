package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dagulv/stock-api/internal/adapter/db"
	"github.com/dagulv/stock-api/internal/env"
	"gopkg.in/yaml.v2"

	migrate "github.com/rubenv/sql-migrate"
)

const (
	configFile    = "dbconfig.yml"
	configEnvType = "development"
)

type configEnv struct {
	DataSource string `yaml:"datasource"`
	Dir        string `yaml:"dir"`
}

var templateContent = `-- +migrate Up

-- +migrate Down
`
var tpl = template.Must(template.New("new_migration").Parse(templateContent))

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := start(ctx, os.Args[1:]); err != nil {
		panic(err)
	}
}

func start(ctx context.Context, args []string) (err error) {
	configEnv, envVars, err := readExpandedConfig()

	db, err := db.Open(ctx, envVars)

	if err != nil {
		return
	}

	if len(args) < 1 {
		return errors.ErrUnsupported
	}

	source := migrate.FileMigrationSource{
		Dir: dir(configEnv.Dir),
	}

	switch args[0] {
	case "up":
		_, err = migrate.ExecMax(db, "postgres", source, migrate.Up, 0)

		if err != nil {
			log.Println(err)
			return errors.ErrUnsupported
		}

		fmt.Println("Applied up migration/s")
	case "down":
		_, err = migrate.ExecMax(db, "postgres", source, migrate.Down, 1)

		if err != nil {
			return err
		}

		fmt.Println("Applied down migration/s")
	case "create":
		if len(args) < 2 {
			return errors.New("a name for the migration is needed")
		}
		if _, err := os.Stat(configEnv.Dir); os.IsNotExist(err) {
			return err
		}

		files, err := os.ReadDir(configEnv.Dir)

		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("%d_%s.sql", len(files)+1, strings.TrimSpace(args[1]))
		pathName := path.Join(configEnv.Dir, fileName)
		f, err := os.Create(pathName)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		if err := tpl.Execute(f, nil); err != nil {
			return err
		}
	}

	return
}

func readExpandedConfig() (config *configEnv, envVars env.Env, err error) {
	if envVars, err = env.GetEnv(dir(".env")); err != nil {
		return
	}

	file, err := os.ReadFile(dir(configFile))
	if err != nil {
		return
	}

	c := make(map[string]*configEnv)
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return
	}

	config = c[configEnvType]

	config.DataSource = os.ExpandEnv(config.DataSource)

	return config, envVars, nil
}

func dir(file string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, file)
}
