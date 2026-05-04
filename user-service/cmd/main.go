package main

import (
	"log"
	"net/http"
	"os"

	"github.com/asdlc-repos/testingnewtaskpage3979/user-service/internal/handlers"
	"github.com/asdlc-repos/testingnewtaskpage3979/user-service/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9091"
	}

	s := store.New()
	h := handlers.New(s)

	mux := http.NewServeMux()
	router := h.RegisterRoutes(mux)

	log.Printf("user-service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
