package main

import (
	"flag"
	"log"
	"net/http"
)

var port = flag.String("port", "", "server port")

func main() {
	flag.Parse()
	cfg := LoadConfig()

	repo := NewInMemoryTodoRepository()
	service := NewTodoService(repo)
	handler := NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/todos", handler.Todos)
	mux.HandleFunc("/todos/", handler.TodoByID)

	listenPort := cfg.Port
	if *port != "" {
		listenPort = *port
	}
	if listenPort == "" {
		listenPort = "8080"
	}

	addr := ":" + listenPort

	db, dbErr := NewDB(cfg)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
