package main

import (
	"context"
	"log"
	"net/http"
	"payment-gateway/internal"
	"payment-gateway/payment/app"
	"payment-gateway/payment/infra"
	"time"
)

func main() {
	ctx, cancel := internal.Context()
	defer cancel()

	mux := http.NewServeMux()
	routes(mux)

	server := http.Server{
		Addr:        ":8080",
		Handler:     mux,
		ReadTimeout: 5 * time.Second,
	}

	log.Println("Starting server on port 8080")

	go func() {
		err := server.ListenAndServe()
		cancel()

		if err != http.ErrServerClosed {
			log.Print(err)
		}
	}()

	<-ctx.Done()
	log.Print("Shutting down server")

	// we give the server some time to shutdown gracefully
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Print("Server gracefully stopped")
}

func routes(mux *http.ServeMux) {
	storage := infra.NewStorage()

	p := app.NewPayment(storage, infra.NewCCProcessor())
	handler := infra.NewHandler(p)

	mux.HandleFunc("/authorize", handler.Authorize)
	mux.HandleFunc("/void", handler.Void)
	mux.HandleFunc("/capture", handler.Capture)
	mux.HandleFunc("/refund", handler.Refund)
}
