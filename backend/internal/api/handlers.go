package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/github/orchestrator/backend/internal/engine"
)

type Handlers struct {
	engine engine.Engine
}

func New(engine engine.Engine) *Handlers {
	return &Handlers{engine: engine}
}

func (h *Handlers) Health(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func (h *Handlers) Ready(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ready"})
}

func (h *Handlers) Live(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "live"})
}

func (h *Handlers) Clusters(c *fiber.Ctx) error {
	clusters, err := h.engine.Clusters()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(clusters)
}

func (h *Handlers) Instances(c *fiber.Ctx) error {
	instances, err := h.engine.Instances()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(instances)
}

func (h *Handlers) Topology(c *fiber.Ctx) error {
	topology, err := h.engine.Topology()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(topology)
}
