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
	userRepo := repository.NewPostgresUserRepository(db)

	todoService := service.NewTodoService(repo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	httpHandler := handler.New(todoService)
	authHandler := handler.NewAuthHandler(authService)

	mux.HandleFunc("/", httpHandler.Home)
	mux.HandleFunc("/health", httpHandler.Health)
	mux.HandleFunc("/todos", httpHandler.Todos)
	mux.HandleFunc("/todos/", httpHandler.TodoByID)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/auth/refresh", authHandler.Refresh)

	rootHandler := middleware.CORS(
		middleware.Recovery(
			middleware.Gzip(mux),
		),
	)

	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, rootHandler))
}
