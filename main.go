package main

import (
	"flag"
	"log"
	"net/http"
)

var port = flag.String("port", "8080", "server port")

func main() {
	flag.Parse()

	repo := NewInMemoryTodoRepository()
	service := NewTodoService(repo)
	handler := NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/todos", handler.Todos)
	mux.HandleFunc("/todos/", handler.TodoByID)

	addr := ":" + *port

	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
