package main

import (
	"context"
	"encoding/json"
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
)

type Cluster struct {
	Name string `json:"name"`
}

type Instance struct {
	Key  string `json:"key"`
	Role string `json:"role"`
}

type Topology struct {
	Clusters  []Cluster  `json:"clusters"`
	Instances []Instance `json:"instances"`
}

type Event struct {
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	logger, err := zap.NewProduction()
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

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ok"})
	})
	app.Get("/readyz", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ready"})
	})
	app.Get("/livez", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "live"})
	})

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Get("/api/v1/clusters", func(c *fiber.Ctx) error {
		clusters := []Cluster{}
		return c.JSON(clusters)
	})
	app.Get("/api/v1/instances", func(c *fiber.Ctx) error {
		instances := []Instance{}
		return c.JSON(instances)
	})
	app.Get("/api/v1/topology", func(c *fiber.Ctx) error {
		topology := Topology{Clusters: []Cluster{}, Instances: []Instance{}}
		return c.JSON(topology)
	})

	app.Get("/events", websocket.New(func(conn *websocket.Conn) {
		defer func() {
			_ = conn.Close()
		}()

		payload, _ := json.Marshal(Event{
			Type:      "hello",
			Message:   "connected",
			Timestamp: time.Now().UTC(),
		})
		_ = conn.WriteMessage(websocket.TextMessage, payload)

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		done := make(chan struct{})
		go func() {
			for range ticker.C {
				message, _ := json.Marshal(Event{
					Type:      "heartbeat",
					Message:   "ok",
					Timestamp: time.Now().UTC(),
				})
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					close(done)
					return
				}
			}
		}()

		for {
			if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
				return
			}
			if _, _, err := conn.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Info("websocket closed", zap.Error(err))
				}
				return
			}
			select {
			case <-done:
				return
			default:
			}
		}
	}))

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("starting api server", zap.String("addr", ":5000"))
		serverErrors <- app.Listen(":5000")
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("failed to shutdown", zap.Error(err))
	}

	logger.Info("shutdown complete")
	fmt.Println("conductor api stopped")
}
