package health

import "github.com/gofiber/fiber/v2"

func Register(app *fiber.App) {
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
}
