package handlers

import (
	"appsku-golang/app/middlewares"
	"context"

	"appsku-golang/app/global-utils/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MainHttpHandler(ctx context.Context) *gin.Engine {
	g := gin.Default()
	g.Use(middlewares.CORSMiddleware(), middlewares.JSONMiddleware(), RequestId(), gin.Recovery())

	return g
}

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get(constants.XRequestId)

		// Create request id with UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Expose it for use in the application
		c.Set(constants.RequestId, requestID)

		c.Writer.Header().Set(constants.XRequestId, requestID)
		c.Next()
	}
}
