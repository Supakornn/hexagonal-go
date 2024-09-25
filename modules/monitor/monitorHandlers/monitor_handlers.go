package monitorHandlers

import (
	"github.com/Supakornn/go-api/config"
	"github.com/Supakornn/go-api/modules/monitor"
	"github.com/gofiber/fiber/v2"
)

type IMonitorHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandlers struct {
	cfg config.IConfig
}

func MonitorHandlers(cfg config.IConfig) IMonitorHandler {
	return &monitorHandlers{
		cfg: cfg,
	}
}

func (h *monitorHandlers) HealthCheck(c *fiber.Ctx) error {
	res := &monitor.Monitor{
		Name:    h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
