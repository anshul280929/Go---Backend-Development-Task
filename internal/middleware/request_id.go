package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestID injects a unique X-Request-Id header into every response
// and stores it in Fiber Locals for downstream use.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Use incoming header if present, otherwise generate a new one.
		requestID := c.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in Locals so other middleware/handlers can access it.
		c.Locals("requestId", requestID)

		// Set on the response.
		c.Set("X-Request-Id", requestID)

		return c.Next()
	}
}
