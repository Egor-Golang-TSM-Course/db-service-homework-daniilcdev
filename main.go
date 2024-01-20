package main

import (
	"database/sql"
	"db-service/handlers/auth"
	"db-service/handlers/comments"
	"db-service/handlers/posts"
	"db-service/handlers/tags"
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

	router := chi.NewRouter()

	userService := auth.NewService().
		WithDb(queries)

	router.Post("/users/register", userService.Register)
	router.Post("/users/login", userService.Login)
	m := middleware.Auth(userService)

	router.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return h
		})

		r.Get("/posts", m.HandlerFunc(posts.GetAllPosts))
		r.Post("/posts", m.HandlerFunc(posts.CreatePost))

		r.Get("/posts/search", m.HandlerFunc(posts.SearchContent))

		r.Get("/posts/{id}", m.HandlerFunc(posts.GetPost))
		r.Put("/posts/{id}", m.HandlerFunc(posts.UpdatePost))
		r.Delete("/posts/{id}", m.HandlerFunc(posts.DeletePost))
	})

	router.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return h
		})

		r.Get("/posts/{postId}/comments", m.HandlerFunc(comments.GetAllComments))
		r.Post("/posts/{postId}/comments", m.HandlerFunc(comments.CreateComment))
	})

	router.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return h
		})

		r.Get("/tags", m.HandlerFunc(tags.GetTags))
		r.Post("/posts/{postId}/tags", m.HandlerFunc(tags.AddTag))
	})

	log.Printf("Starting server on port %s\n", serverPort)

	v1 := chi.NewRouter()
	v1.Mount("/v1", router)
	err = http.ListenAndServe(":"+serverPort, v1)

	if err != nil {
		fmt.Println(err)
	}
}
