package observability

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		requestID, _ := c.Locals("request_id").(string)
		log.Printf("request_id=%s method=%s path=%s status=%d duration=%s", requestID, c.Method(), c.Path(), c.Response().StatusCode(), duration)
		return err
	}
}
