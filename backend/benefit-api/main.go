package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carlosEA28/benefitGuide/internal/auth"
	"github.com/carlosEA28/benefitGuide/internal/catalog"
	"github.com/carlosEA28/benefitGuide/internal/consultations"
	"github.com/carlosEA28/benefitGuide/internal/documents"
	"github.com/carlosEA28/benefitGuide/internal/shared/config"
	"github.com/carlosEA28/benefitGuide/internal/shared/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		slog.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	database, err := db.New(&cfg.Database)
	if err != nil {
		return err
	}
	pool, err := database.DB()
	if err != nil {
		return fmt.Errorf("get PostgreSQL pool: %w", err)
	}
	defer pool.Close()
	pool.SetMaxOpenConns(25)
	pool.SetMaxIdleConns(10)
	pool.SetConnMaxLifetime(30 * time.Minute)
	pingCtx, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	err = pool.PingContext(pingCtx)
	cancelPing()
	if err != nil {
		return fmt.Errorf("ping PostgreSQL: %w", err)
	}

	// MODULOS
	authModule := auth.New(database)
	catalogModule := catalog.New(database)
	consultationsModule := consultations.New(database, authModule.Service, catalogModule.Service)
	documentsModule := documents.New(database, consultationsModule.Service)

	// ROUTER
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer, cors)
	router.Get("/health", healthCheck(pool))
	authModule.RegisterRoutes(router)
	catalogModule.RegisterRoutes(router)
	consultationsModule.RegisterRoutes(router)
	documentsModule.RegisterRoutes(router)

	httpServer := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	serveErrors := make(chan error, 1)
	go func() { serveErrors <- httpServer.ListenAndServe() }()
	slog.Info("HTTP server started", "port", cfg.ServerPort)

	select {
	case err := <-serveErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil
	case <-ctx.Done():
		slog.Info("shutting down HTTP server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		return nil
	}
}

func healthCheck(pool *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := pool.PingContext(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"unavailable"}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
