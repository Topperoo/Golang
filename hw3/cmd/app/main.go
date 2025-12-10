package main

import (
	"homework3/internal/delivery/http"
	"homework3/internal/repository/memory"
	"homework3/internal/usecase"
	"log"
	nethttp "net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	repo := memory.NewBalanceMemoryRepository()
	uc := usecase.NewBalanceUseCase(repo)
	handler := http.NewBalanceHandler(uc)
	router := http.SetupRouter(handler)

	log.Printf("Starting server on port %s", port)
	if err := nethttp.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
