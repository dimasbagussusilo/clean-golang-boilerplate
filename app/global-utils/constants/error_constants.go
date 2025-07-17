package constants

const (
	ERROR_ACTION_NAME_CREATE = "create"
	ERROR_ACTION_NAME_GET    = "get"
	ERROR_ACTION_NAME_UPDATE = "update"
	ERROR_ACTION_NAME_DELETE = "delete"
	ERROR_ACTION_NAME_UPLOAD = "upload"

	ERROR_INVALID_PROCESS                   = "Invalid Process"
	ERROR_BAD_REQUEST_INT_ID_PARAMS         = "Parameter 'id' harus bernilai integer"
	ERROR_DATA_NOT_FOUND                    = "data not found"
	ERROR_INTERNAL_SERVER_1                 = "Ada kesalahan, silahkan coba lagi nanti"
	ERROR_SALESMAN_MOBILE_NUMBER_USED       = "nomor handphone sudah terdaftar sebagai salesman!"
	ERROR_SALESMAN_MOBILE_NUMBER_USED_MSG   = "Maaf, nomor handphone telah terdaftar di Salesman. Silakan gunakan nomor lain!"
	ERROR_SALESMAN_STORE_MOBILE_NUMBER_USED = "Nomor salesman terindikasi sama dengan Nomor user toko, silahkan diperiksa kembali atau gunakan nomor handphone lain"
	ERROR_STORE_MOBILE_NUMBER_USED          = "nomor handphone sudah terdaftar sebagai store!"
	ERROR_STORE_MOBILE_NUMBER_USED_MSG      = "Maaf, nomor handphone telah terdaftar di Toko. Silakan gunakan nomor lain!"

	ERROR_GRPC_INIT_DIAL_FAILED = "Grpc Initial Dial failed"
	ERROR_GRPC_DIAL_FAILED      = "Grpc Dial failed Call %v Service"

	ERROR_SSL_READ_FILE_FAILED = "Failed Read SSL file"
)
