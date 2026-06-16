package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/user/user-management-api/internal/logger"
)

// RequestLogger logs every HTTP request's method, path, status code,
// and duration using Uber Zap. Includes the requestId if available.
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		duration := time.Since(start)
		log := logger.Get()

		// Retrieve requestId set by RequestID middleware.
		requestID, _ := c.Locals("requestId").(string)

		log.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("requestId", requestID),
		)

		return err
	}
}
