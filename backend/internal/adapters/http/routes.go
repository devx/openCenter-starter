package http

import "github.com/gofiber/fiber/v2"

func registerRoutes(app *fiber.App, clusters *ClusterHandler) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/status", func(c *fiber.Ctx) error {
		return WriteJSON(c, fiber.StatusOK, NewSuccess(RequestIDFromCtx(c), fiber.Map{
			"service": "openCenter-base",
			"status":  "ok",
		}))
	})

	v1.Get("/clusters", clusters.listClusters)
	v1.Get("/clusters/:id", clusters.getCluster)
	v1.Post("/clusters", clusters.createCluster)
	v1.Patch("/clusters/:id", clusters.updateCluster)
	v1.Delete("/clusters/:id", clusters.archiveCluster)
}
