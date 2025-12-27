package http

import "github.com/gofiber/fiber/v2"

func registerRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/status", func(c *fiber.Ctx) error {
		return WriteJSON(c, fiber.StatusOK, NewSuccess(RequestIDFromCtx(c), fiber.Map{
			"service": "openCenter-base",
			"status":  "ok",
		}))
	})

	v1.Get("/clusters", listClusters)
	v1.Get("/clusters/:id", getCluster)
	v1.Post("/clusters", createCluster)
	v1.Patch("/clusters/:id", updateCluster)
	v1.Delete("/clusters/:id", archiveCluster)
}
