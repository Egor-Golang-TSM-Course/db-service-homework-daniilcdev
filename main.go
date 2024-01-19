package main

import (
	"context"
	"database/sql"
	"db-service/internal"
	"db-service/internal/database"
	"db-service/middleware"
	"db-service/userService"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	serverPort := os.Getenv("SERVER_PORT")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	queries := database.New(db)

	r := chi.NewRouter()

	userService := userService.NewService().
		WithDb(queries)

	r.Post("/users/register", userService.Register)
	r.Post("/users/login", userService.Login)

	m := middleware.Auth(userService)
	r.Get("/posts", m.HandlerFunc(func(w http.ResponseWriter, r *http.Request, ctx context.Context) {
		dbUser := ctx.Value(middleware.UserData).(*database.User)

		log.Printf("secure /posts, user '%s'\n", dbUser.Name)
		internal.RespondWithJSON(w, 200, dbUser)
	}))

	log.Printf("Starting server on port %s\n", serverPort)
	err = http.ListenAndServe(":"+serverPort, r)

	if err != nil {
		fmt.Println(err)
	}
}
