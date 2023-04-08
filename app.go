package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"tutorial.sqlc.dev/app/tutorial"

	"github.com/apex/gateway/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"

	log "golang.org/x/exp/slog"
)

var GoVersion = runtime.Version()

//go:embed static
var static embed.FS

type Server struct {
	router *chi.Mux
	ctx    context.Context
	db     *tutorial.Queries
}

func NewServer() (*Server, error) {
	srv := Server{
		router: chi.NewRouter(),
		ctx:    context.Background(),
	}

	srv.router.Use(middleware.Logger)

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Error("error connecting to database", err)
	}

	srv.db = tutorial.New(db)

	srv.router.Get("/", srv.handleIndex)
	srv.router.Delete("/author/{id}", srv.handleDeleteAuthor)

	return &srv, nil

}

func main() {

	server, err := NewServer()
	if err != nil {
		log.Error("failed to create server", err)
	}

	if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		log.SetDefault(log.New(log.NewJSONHandler(os.Stdout)))
		err = gateway.ListenAndServe("", server.router)
	} else {
		log.SetDefault(log.New(log.NewTextHandler(os.Stdout)))
		log.Info("local development", "port", os.Getenv("PORT"))
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), server.router)
	}
	log.Error("error listening", err)

}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Version", fmt.Sprintf("%s %s", os.Getenv("version"), GoVersion))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.Must(template.New("").ParseFS(static, "static/index.html"))

	authors, err := s.db.ListAuthors(s.ctx)
	if err != nil {
		log.Error("error listing authors", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("authors", authors)

	err = t.ExecuteTemplate(w, "index.html", struct {
		Authors []tutorial.Author
	}{
		Authors: authors,
	})

	if err != nil {
		log.Error("template failed to parse", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handle delete on /authors/{id}
func (s *Server) handleDeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	log.Info("deleting author", "id", id)
	err := s.db.DeleteAuthor(s.ctx, int64(id))
	if err != nil {
		log.Error("error deleting author", err)
	}
	w.WriteHeader(http.StatusOK)
}
