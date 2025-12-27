package http

import (
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClusterSummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ClusterCreateRequest struct {
	Name string `json:"name"`
}

type ClusterUpdateRequest struct {
	Name   *string `json:"name"`
	Status *string `json:"status"`
}

type clusterStore struct {
	mu       sync.RWMutex
	clusters map[string]ClusterSummary
}

var clusters = &clusterStore{
	clusters: map[string]ClusterSummary{},
}

func listClusters(c *fiber.Ctx) error {
	clusters.mu.RLock()
	defer clusters.mu.RUnlock()

	result := make([]ClusterSummary, 0, len(clusters.clusters))
	for _, cluster := range clusters.clusters {
		result = append(result, cluster)
	}

	return WriteJSON(c, fiber.StatusOK, NewSuccess(RequestIDFromCtx(c), result))
}

func getCluster(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "cluster id is required"))
	}

	clusters.mu.RLock()
	cluster, ok := clusters.clusters[id]
	clusters.mu.RUnlock()

	if !ok {
		return WriteJSON(c, fiber.StatusNotFound, NewError(RequestIDFromCtx(c), "not_found", "cluster not found"))
	}

	return WriteJSON(c, fiber.StatusOK, NewSuccess(RequestIDFromCtx(c), cluster))
}

func createCluster(c *fiber.Ctx) error {
	var payload ClusterCreateRequest
	if err := c.BodyParser(&payload); err != nil {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "invalid JSON payload"))
	}

	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "name is required"))
	}

	cluster := ClusterSummary{
		ID:     uuid.NewString(),
		Name:   name,
		Status: "provisioning",
	}

	clusters.mu.Lock()
	clusters.clusters[cluster.ID] = cluster
	clusters.mu.Unlock()

	return WriteJSON(c, fiber.StatusCreated, NewSuccess(RequestIDFromCtx(c), cluster))
}

func updateCluster(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "cluster id is required"))
	}

	var payload ClusterUpdateRequest
	if err := c.BodyParser(&payload); err != nil {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "invalid JSON payload"))
	}

	clusters.mu.Lock()
	defer clusters.mu.Unlock()

	cluster, ok := clusters.clusters[id]
	if !ok {
		return WriteJSON(c, fiber.StatusNotFound, NewError(RequestIDFromCtx(c), "not_found", "cluster not found"))
	}

	if payload.Name != nil {
		name := strings.TrimSpace(*payload.Name)
		if name == "" {
			return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "name cannot be empty"))
		}
		cluster.Name = name
	}

	if payload.Status != nil {
		status := strings.TrimSpace(*payload.Status)
		if status == "" {
			return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "status cannot be empty"))
		}
		cluster.Status = status
	}

	clusters.clusters[id] = cluster

	return WriteJSON(c, fiber.StatusOK, NewSuccess(RequestIDFromCtx(c), cluster))
}

func archiveCluster(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "cluster id is required"))
	}

	clusters.mu.Lock()
	defer clusters.mu.Unlock()

	cluster, ok := clusters.clusters[id]
	if !ok {
		return WriteJSON(c, fiber.StatusNotFound, NewError(RequestIDFromCtx(c), "not_found", "cluster not found"))
	}

	cluster.Status = "archived"
	clusters.clusters[id] = cluster

	return WriteJSON(c, fiber.StatusOK, NewSuccess(RequestIDFromCtx(c), cluster))
}
