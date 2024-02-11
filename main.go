package main

import (
	"database/sql"
	"db-service/adapters"
	"db-service/handlers"
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
	userService := auth.NewService().WithDb(queries)

	router := chi.NewRouter()
	router.Post("/users/register", handlers.Wrap(userService.Register))
	router.Post("/users/login", handlers.Wrap(userService.Login))
	m := auth.NewMiddleware(userService)

	router.Group(func(r chi.Router) {
		postsStorage := posts.NewStorage(queries)

		r.Get("/posts", handlers.Wrap(postsStorage.GetAllPosts))
		r.Post("/posts", m.HandlerFunc(handlers.WrapCtx(postsStorage.CreatePost)))

		r.Get("/posts/search", m.HandlerFunc(internal.NotImplemented))

		r.Get("/posts/{postId}", handlers.Wrap(postsStorage.GetPost))
		r.Put("/posts/{postId}", m.HandlerFunc(handlers.WrapCtx(postsStorage.UpdatePost)))
		r.Delete("/posts/{postId}", m.HandlerFunc(handlers.WrapCtx(postsStorage.DeletePost)))
	})

	router.Group(func(r chi.Router) {
		commentsStorage := comments.NewStorage(queries)

		r.Get("/posts/{postId}/comments", handlers.Wrap(commentsStorage.GetAllComments))
		r.Post("/posts/{postId}/comments", m.HandlerFunc(handlers.WrapCtx(commentsStorage.CreateComment)))
	})

	router.Group(func(r chi.Router) {
		tqa := adapters.NewTagsQueriesAdapter(queries).WithDb(db)
		tagsStorage := tags.NewStorage(tqa)

		r.Get("/tags", handlers.Wrap(tagsStorage.GetTags))
		r.Post("/posts/{postId}/tags", m.HandlerFunc(handlers.WrapCtx(tagsStorage.AddTag)))
	})

	log.Printf("Starting server on port %s\n", serverPort)

	v1 := chi.NewRouter()
	v1.Mount("/v1", router)

	if err = http.ListenAndServe(":"+serverPort, v1); err != nil {
		fmt.Println(err)
	}
}
