package main

import (
	"context"
	"log"

	"github.com/dagulv/stock-api/internal/api/server"
	"github.com/dagulv/stock-api/internal/db"
	"github.com/dagulv/stock-api/internal/models"
	"github.com/dagulv/stock-api/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbP *pgxpool.Pool
	v   *validator.Validate
)

func main() {
	ctx := context.Background()

	Start(ctx)
}

func Start(ctx context.Context) (err error) {
	log.Println("Connecting to database...")
	dbP, err = db.Connect(ctx)

	log.Println("db conn", dbP)
	if err != nil {
		return
	}

	defer dbP.Close()

	v = validator.New()

	userService := services.UserService{
		Store:    db.UserStore(dbP),
		Validate: v,
	}

	s := server.Server{
		UserService: &userService,
	}

	if err = createInitialUser(ctx, &userService); err != nil {
		return
	}

	return s.Start(ctx)
}

func createInitialUser(ctx context.Context, userService *services.UserService) (err error) {
	email := "admin@test.com"

	var user models.User

	if err = userService.GetByEmail(ctx, email, &user); err == nil {
		return
	}

	user.Active = true
	user.Name = "admin"
	user.Email = email
	password := "admin"

	if err = userService.Create(ctx, &user); err != nil {
		return
	}

	err = userService.SetPassword(ctx, password, user.Id)

	return
}
