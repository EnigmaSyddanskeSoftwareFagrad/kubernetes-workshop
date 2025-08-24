package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr: "0.0.0.0:80",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Server Control</title>
			</head>
			<body>
				<form action="/shutdown" method="post">
					<button type="submit">Shutdown Server</button>
				</form>
			</body>
			</html>`
		fmt.Fprint(w, html)
	})

	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Shutdown command received")
		shutdownChan <- os.Interrupt
	})

	go func() {
		fmt.Printf("server started at %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	<-shutdownChan
	fmt.Println("Shutdown signal received, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Graceful shutdown failed: %s\n", err)
	} else {
		fmt.Println("Server shut down gracefully.")
	}
}
