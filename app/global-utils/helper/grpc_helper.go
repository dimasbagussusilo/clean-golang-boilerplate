package helper

import (
	"appsku-golang/app/global-utils/constants"
	"fmt"
	"time"
)

// colorForStatus determines the color based on the status code
func colorForStatus(code int32) string {
	switch {
	case code >= 200 && code < 300:
		return constants.GreenColor
	case code >= 300 && code < 400:
		return constants.YellowColor
	case code >= 400 && code < 500:
		return constants.RedColor
	case code >= 500:
		return constants.MagentaColor
	default:
		return constants.BlueColor
	}
}

// colorForStatus determines the color based on the status code
func colorBgForStatus(code int32) string {
	switch {
	case code >= 200 && code < 300:
		return constants.GreenBg
	case code >= 300 && code < 400:
		return constants.YellowBg
	case code >= 400 && code < 500:
		return constants.RedBg
	case code >= 500:
		return constants.MagentaBg
	default:
		return constants.BlueBg
	}
}

func GrpcLogger(path string, latency time.Time, statusCode int32) {
	// Get status code and assign color
	statusColor := colorBgForStatus(statusCode)
	now := time.Now()
	fmt.Printf("[GRPC] %s |%s %3d %s| %13v | %s\n",
		now.Format("2006/01/02 - 15:04:05"), // Waktu
		statusColor,                         // Set color background start
		statusCode,                          // Status code
		constants.ResetColor,                // Set color background stop
		time.Since(latency),                 // Durasi request
		path,                                // Path endpoint
	)
}
