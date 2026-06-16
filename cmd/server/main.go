package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/user/user-management-api/config"
	"github.com/user/user-management-api/db/sqlc"
	"github.com/user/user-management-api/internal/handler"
	"github.com/user/user-management-api/internal/logger"
	"github.com/user/user-management-api/internal/middleware"
	"github.com/user/user-management-api/internal/repository"
	"github.com/user/user-management-api/internal/routes"
	"github.com/user/user-management-api/internal/service"
)

func main() {
	// ── Logger ───────────────────────────────────────────────
	logger.Init()
	defer logger.Sync()
	log := logger.Get()

	// ── Configuration ────────────────────────────────────────
	cfg := config.Load()

	// ── Database ─────────────────────────────────────────────
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("failed to ping database", zap.Error(err))
	}
	log.Info("connected to database", zap.String("host", cfg.DBHost), zap.String("port", cfg.DBPort))

	// ── Dependency Wiring ────────────────────────────────────
	queries := sqlc.New(pool)
	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// ── Fiber App ────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName: "User Management API v1.0.0",
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Register routes
	routes.Setup(app, userHandler)

	// ── Graceful Shutdown ────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.ServerPort)
		log.Info("starting server", zap.String("address", addr))
		if err := app.Listen(addr); err != nil {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	<-quit
	log.Info("shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Error("server forced shutdown", zap.Error(err))
	}

	log.Info("server stopped gracefully")
}
