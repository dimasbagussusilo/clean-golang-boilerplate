package constants

// ANSI color codes
const (
	ResetColor = "\033[0m"

	GreenColor   = "\033[32m" // 2xx - Success
	YellowColor  = "\033[33m" // 3xx - Redirect
	BlueColor    = "\033[34m" // 1xx - Informational
	RedColor     = "\033[31m" // 4xx - Client Error
	MagentaColor = "\033[35m" // 5xx - Server Error

	GreenBg   = "\033[42m\033[97m" // Green background, white text (2xx)
	YellowBg  = "\033[43m\033[97m" // Yellow background, white text (3xx)
	RedBg     = "\033[41m\033[97m" // Red background, white text (4xx)
	MagentaBg = "\033[45m\033[97m" // Magenta background, white text (5xx)
	BlueBg    = "\033[44m\033[97m" // Blue background, white text (1xx)
)
