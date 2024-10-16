package main

import (
	"ItDevTest/internal/handler"
	"ItDevTest/internal/repository"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "merey"
	password = "postgres"
	dbname   = "it_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Connected to database")

	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/createUser", userHandler.CreateUserHandler)
	http.ListenAndServe(":8085", r)
}
