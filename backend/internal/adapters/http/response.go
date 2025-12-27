package http

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

func NewError(code, message string) ErrorResponse {
	return ErrorResponse{Error: ErrorBody{Code: code, Message: message}}
}

func NewSuccess[T any](data T) SuccessResponse[T] {
	return SuccessResponse[T]{Data: data}
}

func WriteJSON(c *fiber.Ctx, status int, payload any) error {
	return c.Status(status).JSON(payload)
}
