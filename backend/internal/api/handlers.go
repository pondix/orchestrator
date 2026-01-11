package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
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

type Handlers struct{}

func New() *Handlers {
	return &Handlers{}
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
	clusters := []Cluster{}
	return c.JSON(clusters)
}

func (h *Handlers) Instances(c *fiber.Ctx) error {
	instances := []Instance{}
	return c.JSON(instances)
}

func (h *Handlers) Topology(c *fiber.Ctx) error {
	topology := Topology{Clusters: []Cluster{}, Instances: []Instance{}}
	return c.JSON(topology)
}
