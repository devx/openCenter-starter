package http

import "github.com/gofiber/fiber/v2"

type ClusterSummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func listClusters(c *fiber.Ctx) error {
	clusters := []ClusterSummary{}
	return WriteJSON(c, fiber.StatusOK, NewSuccess(RequestIDFromCtx(c), clusters))
}
