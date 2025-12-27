package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/devx/openCenter-starter/backend/internal/adapters/memory"
	"github.com/devx/openCenter-starter/backend/internal/adapters/postgres"
	"github.com/devx/openCenter-starter/backend/internal/config"
	"github.com/devx/openCenter-starter/backend/internal/health"
	"github.com/devx/openCenter-starter/backend/internal/observability"
	"github.com/devx/openCenter-starter/backend/internal/ports"
)

type Server struct {
	app *fiber.App
}

func New(cfg config.Config) (*Server, error) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			code := "internal_error"
			message := "internal server error"

			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				status = fiberErr.Code
				code = http.StatusText(fiberErr.Code)
				if fiberErr.Message != "" {
					message = fiberErr.Message
				}
			}

			if code == "" {
				code = "internal_error"
			}

			return WriteJSON(c, status, NewError(RequestIDFromCtx(c), code, message))
		},
	})
	app.Use(observability.RequestID())
	app.Use(observability.RequestLogger())
	health.Register(app)
	clusterStore, err := resolveClusterStore(cfg)
	if err != nil {
		return nil, err
	}
	clusterHandler := NewClusterHandler(clusterStore)
	registerRoutes(app, clusterHandler)

	return &Server{app: app}, nil
}

func resolveClusterStore(cfg config.Config) (ports.ClusterStore, error) {
	if cfg.DatabaseURL == "" {
		return memory.NewClusterStore(), nil
	}

	return postgres.NewClusterStore(context.Background(), cfg.DatabaseURL)
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}
