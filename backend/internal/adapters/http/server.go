package http

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/devx/openCenter-starter/backend/internal/health"
	"github.com/devx/openCenter-starter/backend/internal/observability"
)

type Server struct {
	app *fiber.App
}

func New() *Server {
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

			return WriteJSON(c, status, NewError(code, message))
		},
	})
	app.Use(observability.RequestLogger())
	health.Register(app)
	registerRoutes(app)

	return &Server{app: app}
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}
