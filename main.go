package main

import (
	"database/sql"
	"db-service/handlers/auth"
	"db-service/handlers/comments"
	"db-service/handlers/posts"
	"db-service/handlers/tags"
	"db-service/internal"
	"db-service/internal/database"
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
	m := auth.NewMiddleware(userService)

	router.Group(func(r chi.Router) {
		postsStorage := posts.NewStorage(queries)

		r.Get("/posts", postsStorage.GetAllPosts)
		r.Post("/posts", m.HandlerFunc(postsStorage.CreatePost))

		r.Get("/posts/search", m.HandlerFunc(internal.NotImplemented))

		r.Get("/posts/{postId}", postsStorage.GetPost)
		r.Put("/posts/{postId}", m.HandlerFunc(postsStorage.UpdatePost))
		r.Delete("/posts/{postId}", m.HandlerFunc(postsStorage.DeletePost))
	})

	router.Group(func(r chi.Router) {
		commentsStorage := comments.NewStorage(queries)

		r.Get("/posts/{postId}/comments", commentsStorage.GetAllComments)
		r.Post("/posts/{postId}/comments", m.HandlerFunc(commentsStorage.CreateComment))
	})

	router.Group(func(r chi.Router) {
		r.Get("/tags", m.HandlerFunc(tags.GetTags))
		r.Post("/posts/{postId}/tags", m.HandlerFunc(internal.NotImplemented))
	})

	log.Printf("Starting server on port %s\n", serverPort)

	v1 := chi.NewRouter()
	v1.Mount("/v1", router)
	err = http.ListenAndServe(":"+serverPort, v1)

	if err != nil {
		fmt.Println(err)
	}
}
