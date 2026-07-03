package main

import (
	"context"
	"embed"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-uptime-monitor/cli"
	"go-uptime-monitor/config"
	"go-uptime-monitor/database"
	"go-uptime-monitor/handlers"
	"go-uptime-monitor/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func loadEnv() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Warn().Err(err).Msg("Warning: failed to load .env file")
	}
}

func setupLogger(level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	
	// Convert level string to zerolog.Level
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		l = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(l)

	// If we're not running as a server (or just for nice local output), we could use ConsoleWriter.
	// But the user requested structured JSON logging for production, which is the default for zerolog.
	// We'll leave it as JSON.
}

func main() {
	// Initialize a basic logger before config is loaded
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	loadEnv()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	setupLogger(cfg.LogLevel)

	if len(os.Args) < 2 {
		cli.PrintUsage()
		return
	}

	switch os.Args[1] {
	case "s", "server":
		runServer(cfg)
	case "migrate":
		database.RunMigrations(embedMigrations, cfg.Database)
	default:
		db := database.Connect(cfg.Database)
		cli.Run(db, os.Args[1:])
	}
}

func runServer(cfg *config.Config) {
	gin.SetMode(cfg.GinMode)

	// Use gin.New() instead of Default() to exclude the default Logger
	r := gin.New()
	
	// Add our structured logger and standard Recovery middleware
	r.Use(middlewares.ZerologLogger())
	r.Use(gin.Recovery())

	gormDB := database.Connect(cfg.Database)
	h := handlers.NewHandler(gormDB)

	r.GET("/health", h.Health)

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!DOCTYPE html><html><head><title>Hello</title></head><body><h1>Hello, World!</h1></body></html>`))
	})

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.HTTPPort).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
