package middleware

import (
	"appsku-golang/app/global-utils/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware memastikan setiap request memiliki Request ID
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Periksa apakah Request ID sudah ada dalam header
		requestID := c.GetHeader(constants.XRequestId)
		if requestID == "" {
			// Jika tidak ada, buat Request ID baru
			requestID = uuid.New().String()
			c.Request.Header.Set(constants.XRequestId, requestID)
		}

		// Simpan Request ID ke dalam context untuk penggunaan lebih lanjut
		c.Set(constants.RequestId, requestID)

		// Tambahkan Request ID ke header respons
		c.Writer.Header().Set(constants.XRequestId, requestID)

		// Lanjutkan ke handler berikutnya
		c.Next()
	}
}
