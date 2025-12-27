package http

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/devx/openCenter-starter/backend/internal/ports"
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

type ClusterHandler struct {
	store ports.ClusterStore
}

func NewClusterHandler(store ports.ClusterStore) *ClusterHandler {
	return &ClusterHandler{store: store}
}

func (h *ClusterHandler) listClusters(c *fiber.Ctx) error {
	statusFilter := strings.TrimSpace(c.Query("status"))
	namePrefix := strings.TrimSpace(c.Query("name_prefix"))
	idPrefix := strings.TrimSpace(c.Query("id_prefix"))
	limit := parseQueryInt(c, "limit", 50)
	offset := parseQueryInt(c, "offset", 0)
	if limit < 1 {
		limit = 1
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	filter := ports.ClusterFilter{
		Status:     statusFilter,
		NamePrefix: namePrefix,
		IDPrefix:   idPrefix,
		Limit:      limit,
		Offset:     offset,
	}

	items, total, err := h.store.List(c.Context(), filter)
	if err != nil {
		return WriteJSON(c, fiber.StatusInternalServerError, NewError(RequestIDFromCtx(c), "internal_error", "unable to list clusters"))
	}

	result := make([]ClusterSummary, 0, len(items))
	for _, cluster := range items {
		result = append(result, ClusterSummary{
			ID:     cluster.ID,
			Name:   cluster.Name,
			Status: cluster.Status,
		})
	}

	if offset >= total {
		return WriteJSON(c, fiber.StatusOK, NewSuccessWithPagination(
			RequestIDFromCtx(c),
			[]ClusterSummary{},
			PaginationMeta{Total: total, Limit: limit, Offset: offset},
		))
	}

	return WriteJSON(c, fiber.StatusOK, NewSuccessWithPagination(
		RequestIDFromCtx(c),
		result,
		PaginationMeta{Total: total, Limit: limit, Offset: offset},
	))
}

func (h *ClusterHandler) getCluster(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "cluster id is required"))
	}

	cluster, ok, err := h.store.Get(c.Context(), id)
	if err != nil {
		return WriteJSON(c, fiber.StatusInternalServerError, NewError(RequestIDFromCtx(c), "internal_error", "unable to read cluster"))
	}

	if !ok {
		return WriteJSON(c, fiber.StatusNotFound, NewError(RequestIDFromCtx(c), "not_found", "cluster not found"))
	}

	return WriteJSON(c, fiber.StatusOK, NewSuccess(RequestIDFromCtx(c), ClusterSummary{
		ID:     cluster.ID,
		Name:   cluster.Name,
		Status: cluster.Status,
	}))
}

func (h *ClusterHandler) createCluster(c *fiber.Ctx) error {
	var payload ClusterCreateRequest
	if err := c.BodyParser(&payload); err != nil {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "invalid JSON payload"))
	}

	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "name is required"))
	}

	cluster, err := h.store.Create(c.Context(), name)
	if err != nil {
		return WriteJSON(c, fiber.StatusInternalServerError, NewError(RequestIDFromCtx(c), "internal_error", "unable to create cluster"))
	}

	return WriteJSON(c, fiber.StatusCreated, NewSuccess(RequestIDFromCtx(c), ClusterSummary{
		ID:     cluster.ID,
		Name:   cluster.Name,
		Status: cluster.Status,
	}))
}

func (h *ClusterHandler) updateCluster(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "cluster id is required"))
	}

	var payload ClusterUpdateRequest
	if err := c.BodyParser(&payload); err != nil {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "invalid JSON payload"))
	}

	if payload.Name != nil {
		name := strings.TrimSpace(*payload.Name)
		if name == "" {
			return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "name cannot be empty"))
		}
		payload.Name = &name
	}

	if payload.Status != nil {
		status := strings.TrimSpace(*payload.Status)
		if status == "" {
			return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "status cannot be empty"))
		}
		payload.Status = &status
	}

	update := ports.ClusterUpdate{
		Name:   payload.Name,
		Status: payload.Status,
	}

	cluster, ok, err := h.store.Update(c.Context(), id, update)
	if err != nil {
		return WriteJSON(c, fiber.StatusInternalServerError, NewError(RequestIDFromCtx(c), "internal_error", "unable to update cluster"))
	}
	if !ok {
		return WriteJSON(c, fiber.StatusNotFound, NewError(RequestIDFromCtx(c), "not_found", "cluster not found"))
	}

	return WriteJSON(c, fiber.StatusOK, NewSuccess(RequestIDFromCtx(c), ClusterSummary{
		ID:     cluster.ID,
		Name:   cluster.Name,
		Status: cluster.Status,
	}))
}

func (h *ClusterHandler) archiveCluster(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return WriteJSON(c, fiber.StatusBadRequest, NewError(RequestIDFromCtx(c), "invalid_request", "cluster id is required"))
	}

	cluster, ok, err := h.store.Archive(c.Context(), id)
	if err != nil {
		return WriteJSON(c, fiber.StatusInternalServerError, NewError(RequestIDFromCtx(c), "internal_error", "unable to archive cluster"))
	}
	if !ok {
		return WriteJSON(c, fiber.StatusNotFound, NewError(RequestIDFromCtx(c), "not_found", "cluster not found"))
	}

	return WriteJSON(c, fiber.StatusOK, NewSuccess(RequestIDFromCtx(c), ClusterSummary{
		ID:     cluster.ID,
		Name:   cluster.Name,
		Status: cluster.Status,
	}))
}

func parseQueryInt(c *fiber.Ctx, key string, fallback int) int {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
