package middlewares

import (
	"appsku-golang/app/config"
	"appsku-golang/app/global-utils/constants"
	"appsku-golang/app/global-utils/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

// BasicAuthMiddleware returns a Gin middleware function that implements basic authentication
func BasicAuthMiddleware() gin.HandlerFunc {
	cfg := config.Get()
	username := cfg.AuthBasic.Username
	password := cfg.AuthBasic.Password

	return func(c *gin.Context) {
		auth := c.GetHeader(constants.Authorization)
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Unauthorized: Basic authentication required",
			})
			return
		}

		if !strings.HasPrefix(auth, "Basic ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Unauthorized: Invalid authentication method",
			})
			return
		}

		expectedAuthBase64 := "Basic " + helper.BasicAuth(username, password)

		if auth != expectedAuthBase64 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Unauthorized: Invalid credentials",
			})
			return
		}

		c.Next()
	}
}

// FiberBasicAuthMiddleware returns a Fiber middleware function that implements basic authentication
func FiberBasicAuthMiddleware() fiber.Handler {
	cfg := config.Get()
	username := cfg.AuthBasic.Username
	password := cfg.AuthBasic.Password

	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		auth := c.Get(constants.Authorization)
		if auth == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized: Basic authentication required",
			})
		}

		// Check if it's a Basic auth
		if !strings.HasPrefix(auth, "Basic ") {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized: Invalid authentication method",
			})
		}

		// Compare with expected credentials
		expectedAuthBase64 := "Basic " + helper.BasicAuth(username, password)

		if auth != expectedAuthBase64 {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized: Invalid credentials",
			})
		}

		// Authentication successful, continue
		return c.Next()
	}
}
