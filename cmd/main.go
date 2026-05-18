package main

import (
	"flag"
	"log"
	"net/http"

	"maxsasi/internal/config"
	"maxsasi/internal/database"
	"maxsasi/internal/handler"
	"maxsasi/internal/repository"
	"maxsasi/internal/service"
	"maxsasi/pkg/middleware"
)

func main() {
	port := flag.String("port", "", "server port")
	flag.Parse()

	cfg := config.Load()
	mux := http.NewServeMux()

	listenPort := cfg.Port
	if *port != "" {
		listenPort = *port
	}
	if listenPort == "" {
		listenPort = "8080"
	}

	addr := ":" + listenPort

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresTodoRepository(db)
	todoService := service.NewTodoService(repo)
	httpHandler := handler.New(todoService)

	mux.HandleFunc("/", httpHandler.Home)
	mux.HandleFunc("/health", httpHandler.Health)
	mux.HandleFunc("/todos", httpHandler.Todos)
	mux.HandleFunc("/todos/", httpHandler.TodoByID)

	rootHandler := middleware.CORS(
		middleware.Recovery(
			middleware.Auth(cfg.AuthToken)(
				middleware.Gzip(mux),
			),
		),
	)

	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, rootHandler))
}
