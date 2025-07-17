package handlers

import (
	"appsku-golang/app/middlewares"
	"context"

	"appsku-golang/app/global-utils/constants"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
)

func MainFiberHandler(ctx context.Context) *fiber.App {
	f := fiber.New(
		fiber.Config{
			Prefork:       true,
			CaseSensitive: true,
			//StrictRouting: true,
			ServerHeader: "Fiber",
			AppName:      "App v1.0.1",
		},
	)
	f.Use(middlewares.FiberCORSMiddleware(), middlewares.FiberJSONMiddleware(), FiberRequestId(), recover.New())

	return f
}

func FiberRequestId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(constants.XRequestId)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Locals(constants.RequestId, requestID)

		c.Set(constants.XRequestId, requestID)
		return c.Next()
	}
}
