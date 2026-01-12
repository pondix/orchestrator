package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/websocket/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/github/orchestrator/backend/internal/api"
	"github.com/github/orchestrator/backend/internal/config"
	"github.com/github/orchestrator/backend/internal/engine"
	"github.com/github/orchestrator/backend/internal/events"
	"github.com/github/orchestrator/backend/internal/logging"
)

func main() {
	cfg := config.Load()

	logger, err := logging.New()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	app := fiber.New(fiber.Config{
		AppName: "conductor-api",
	})

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		status := c.Response().StatusCode()
		logger.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("duration", time.Since(start)),
		)
		return err
	})

	hub := events.NewHub()
	engineStub := engine.NewStub()
	handlers := api.New(engineStub)

	app.Get("/healthz", handlers.Health)
	app.Get("/readyz", handlers.Ready)
	app.Get("/livez", handlers.Live)

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Get("/api/v1/clusters", handlers.Clusters)
	app.Get("/api/v1/instances", handlers.Instances)
	app.Get("/api/v1/topology", handlers.Topology)

	app.Get("/events", websocket.New(func(conn *websocket.Conn) {
		defer func() {
			hub.Unregister(conn)
			_ = conn.Close()
		}()

		hub.Register(conn)
		hub.Broadcast(events.Event{
			Type:      "hello",
			Message:   "connected",
			Timestamp: time.Now().UTC(),
		})

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		done := make(chan struct{})
		go func() {
			for range ticker.C {
				hub.Broadcast(events.Event{
					Type:      "heartbeat",
					Message:   "ok",
					Timestamp: time.Now().UTC(),
				})
				select {
				case <-done:
					return
				default:
				}
			}
		}()

		for {
			if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
				close(done)
				return
			}
			if _, _, err := conn.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Info("websocket closed", zap.Error(err))
				}
				close(done)
				return
			}
		}
	}))

	eventsTicker := time.NewTicker(30 * time.Second)
	go func() {
		for range eventsTicker.C {
			hub.Broadcast(events.Event{
				Type:      "status",
				Message:   "idle",
				Timestamp: time.Now().UTC(),
			})
		}
	}()

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("starting api server", zap.String("addr", cfg.APIAddr))
		serverErrors <- app.Listen(cfg.APIAddr)
	}()

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("api server error", zap.Error(err))
		}
	case sig := <-shutdownSignals:
		logger.Info("shutdown signal received", zap.String("signal", sig.String()))
	}

	eventsTicker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("failed to shutdown", zap.Error(err))
	}

	logger.Info("shutdown complete")
	fmt.Println("conductor api stopped")
}
