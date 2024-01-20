package main

import (
	"database/sql"
	"db-service/handlers/auth"
	"db-service/handlers/comments"
	"db-service/handlers/posts"
	"db-service/internal/database"
	"db-service/middleware"
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

	userService := auth.NewService().
		WithDb(queries)

	r.Post("/users/register", userService.Register)
	r.Post("/users/login", userService.Login)
	m := middleware.Auth(userService)

	r.Group(func(router chi.Router) {
		router.Use(func(h http.Handler) http.Handler {
			return h
		})

		router.Post("/posts", m.HandlerFunc(posts.CreatePost))
		router.Get("/posts", m.HandlerFunc(posts.GetAllPosts))

		router.Get("/posts/{id}", m.HandlerFunc(posts.GetPost))
		router.Put("/posts/{id}", m.HandlerFunc(posts.UpdatePost))
		router.Delete("/posts/{id}", m.HandlerFunc(posts.DeletePost))
	})

	r.Group(func(router chi.Router) {
		router.Use(func(h http.Handler) http.Handler {
			return h
		})

		router.Post("/posts/{postId}/comments", m.HandlerFunc(comments.CreateComment))
		router.Get("/posts/{postId}/comments", m.HandlerFunc(comments.GetAllComments))
	})

	log.Printf("Starting server on port %s\n", serverPort)

	v1 := chi.NewRouter()
	v1.Mount("/v1", r)
	err = http.ListenAndServe(":"+serverPort, v1)

	if err != nil {
		fmt.Println(err)
	}
}
