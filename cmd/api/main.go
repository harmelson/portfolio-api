package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/harmelson/tocouaboa-portfolio/internal/config"
	"github.com/harmelson/tocouaboa-portfolio/internal/db"
	httpRouter "github.com/harmelson/tocouaboa-portfolio/internal/http"
	"github.com/harmelson/tocouaboa-portfolio/internal/http/handlers"
	"github.com/harmelson/tocouaboa-portfolio/internal/repository"
	"github.com/harmelson/tocouaboa-portfolio/internal/service"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting server on port %s", cfg.Port)

	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Database connection established")

	userRepo := repository.NewUserRepository(pool)
	planRepo := repository.NewPlanRepository(pool)
	subRepo := repository.NewSubscriptionRepository(pool)

	userService := service.NewUserService(pool, userRepo, planRepo, subRepo)
	subService := service.NewSubscriptionService(subRepo, planRepo)

	userHandler := handlers.NewUserHandler(userService)
	subHandler := handlers.NewSubscriptionHandler(subService)

	h := &httpRouter.Handlers{
		User:         userHandler,
		Subscription: subHandler,
	}

	deps := &httpRouter.Dependencies{
		Handlers:    h,
		UserService: userService,
	}
	router := httpRouter.NewRouter(deps)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server running on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
