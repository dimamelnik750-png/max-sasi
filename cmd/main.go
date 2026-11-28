package main

import (
	"flag"
	"log"
	"net/http"

	"maxsasi/internal/cache"
	"maxsasi/internal/config"
	"maxsasi/internal/database"
	"maxsasi/internal/handler"
	"maxsasi/internal/repository"
	"maxsasi/internal/service"
)

func main() {
	port := flag.String("port", "", "server port")
	flag.Parse()

	cfg := config.Load()

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
	redisCache := cache.NewRedisCache(cfg.RedisHost, cfg.RedisPort)

	todoService := service.NewTodoService(repo)
	cachedTodoService := service.NewCachedTodoService(todoService, redisCache)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	httpHandler := handler.New(cachedTodoService)
	authHandler := handler.NewAuthHandler(authService)
	mux := handler.NewRouter(httpHandler, authHandler)

	rootHandler := handler.CORS(
		handler.Recovery(
			handler.JWTAuth(cfg.JWTSecret)(
				handler.Gzip(mux),
			),
		),
	)

	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, rootHandler))
}
