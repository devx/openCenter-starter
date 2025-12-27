package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/devx/openCenter-starter/backend/internal/health"
	"github.com/devx/openCenter-starter/backend/internal/observability"
)

type Server struct {
	app *fiber.App
}

func New() *Server {
	app := fiber.New()
	app.Use(observability.RequestLogger())
	health.Register(app)

	return &Server{app: app}
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}
