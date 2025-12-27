package http

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error ErrorBody    `json:"error"`
	Meta  ResponseMeta `json:"meta"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse[T any] struct {
	Data T            `json:"data"`
	Meta ResponseMeta `json:"meta"`
}

type ResponseMeta struct {
	RequestID  string          `json:"request_id"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

type PaginationMeta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func NewError(requestID, code, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorBody{Code: code, Message: message},
		Meta:  ResponseMeta{RequestID: requestID},
	}
}

func NewSuccess[T any](requestID string, data T) SuccessResponse[T] {
	return SuccessResponse[T]{
		Data: data,
		Meta: ResponseMeta{RequestID: requestID},
	}
}

func NewSuccessWithPagination[T any](requestID string, data T, pagination PaginationMeta) SuccessResponse[T] {
	return SuccessResponse[T]{
		Data: data,
		Meta: ResponseMeta{RequestID: requestID, Pagination: &pagination},
	}
}

func RequestIDFromCtx(c *fiber.Ctx) string {
	requestID, _ := c.Locals("request_id").(string)
	return requestID
}

func WriteJSON(c *fiber.Ctx, status int, payload any) error {
	return c.Status(status).JSON(payload)
}
