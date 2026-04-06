package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "", "server port")
	flag.Parse()
	cfg := LoadConfig()

	mux := http.NewServeMux()

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

	if err := InitDB(db); err != nil {
		log.Fatal(err)
	}

	repo := NewPostgresTodoRepository(db)
	service := NewTodoService(repo)
	handler := NewHandler(service)

	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/todos", handler.Todos)
	mux.HandleFunc("/todos/", handler.TodoByID)

	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
