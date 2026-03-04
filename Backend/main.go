package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/YimingChen/P4SBU/Backend/config"
	mongodb "github.com/YimingChen/P4SBU/Backend/database"
	"github.com/YimingChen/P4SBU/Backend/middleware"
	"github.com/YimingChen/P4SBU/Backend/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using environment variables")
	}

	// Load configuration
	cfg := config.NewConfig()

	// Initialize MongoDB connection
	mongoClient, err := mongodb.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer mongoClient.Disconnect(ctx)

	// Create repository instance
	repo := mongodb.NewRepository(mongoClient, cfg.MongoDB.Database)

	// Create router and register routes
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.LoggingMiddleware)

	// API routes with versioning
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	apiRouter.Use(middleware.CorsMiddleware)

	// Register API routes
	routes.RegisterRoutes(apiRouter, repo)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s\n", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("Server shutting down...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped")
}
