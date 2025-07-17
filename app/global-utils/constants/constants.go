package constants

import "time"

// http
const (
	XRequestId    = "X-Request-Id"
	RequestId     = "requestId"
	XChannel      = "X-Channel"
	Authorization = "Authorization"
)

// date format constants
const (
	DATE_FORMAT_COMMON             = "2006-01-02"
	DATE_FORMAT_EXPORT             = "20060102-150405"
	DATE_FORMAT_EXPORT_CREATED_AT  = "2006-01-02 15:04:05"
	DATE_FORMAT_CODE_GENERATOR     = "20060102150405"
	DATE_TIME_FORMAT_COMON         = "2006-01-02 15:04:05"
	DATE_TIME_ZERO_HOUR_ADDITIONAL = " 00:00:00"
)

// so do and status constants
const (
	SO_STATUS_APPV   = "APPV"
	SO_STATUS_REAPPV = "REAPPV"
	SO_STATUS_RJC    = "RJC"
	SO_STATUS_CNCL   = "CNCL"
	SO_STATUS_ORDPRT = "ORDPRT"
	SO_STATUS_ORDCLS = "ORDCLS"
	SO_STATUS_CLS    = "CLS"
	SO_STATUS_PEND   = "PEND"
	SO_STATUS_OPEN   = "OPEN"

	DO_STATUS_CANCEL = "SJCNCL"
	DO_STATUS_CLOSED = "SJCLS"
	DO_STATUS_OPEN   = "SJCR"

	ORDER_STATUS_OPEN      = "open"
	ORDER_STATUS_CANCELLED = "cancelled"
	ORDER_STATUS_CLOSED    = "closed"
	ORDER_STATUS_PARTIAL   = "partial"
	ORDER_STATUS_PENDING   = "pending"
	ORDER_STATUS_REJECTED  = "rejected"
)

// file type constants
const (
	FILE_EXCEL_TYPE = "xlsx"
	FILE_CSV_TYPE   = "csv"
)

const (
	DEFAULT_CACHE_TIME = 10 * time.Minute
	DUMMY_EMAIL        = "@dummy.dbo"

	// WhatsappCheckAccountStatus
	WASuccess = "valid"
	WAPending = "processing"
	WAFailed  = "failed"
	WANoValid = "invalid"

	WANoAccount         = "NO_ACCOUNT"
	UserNoWhatsapp      = "no"
	UserHasWhatsapp     = "valid"
	UserWaitingWhatsapp = "waiting"
)

// status file upload
const (
	STATUS_UPLOAD_FILE_OPEN       = "uploaded"
	STATUS_UPLOAD_FILE_PROCESSING = "processing"
	STATUS_UPLOAD_FILE_FAILED     = "failed"
	STATUS_UPLOAD_FILE_SUCCESS    = "success"
	STATUS_UPLOAD_FILE_PENDING    = "pending"
)

// Environment application
const (
	ENV_PRODUCTION  = "production"
	ENV_PRELIVE     = "prelive"
	ENV_STAGING     = "staging"
	ENV_DEVELOPMENT = "development"
	ENV_UNDEFINED   = "undefined_env"
)
