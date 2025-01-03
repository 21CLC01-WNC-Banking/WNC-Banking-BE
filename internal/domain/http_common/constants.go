package httpcommon

type errorResponseCode struct {
	InvalidRequest      string
	InternalServerError string
	RecordNotFound      string
	MissingIdParameter  string
	InvalidDataType     string
	InvalidUserInfo     string
	Unauthorized        string
	Forbidden           string
	TimeoutRequest      string
}

var ErrorResponseCode = errorResponseCode{
	InvalidRequest:      "INVALID_REQUEST",
	InternalServerError: "INTERNAL_SERVER_ERROR",
	RecordNotFound:      "RECORD_NOT_FOUND",
	MissingIdParameter:  "MISSING_ID_PARAMETER",
	InvalidDataType:     "INVALID_DATA_TYPE",
	InvalidUserInfo:     "INVALID_USER_INFO",
	Unauthorized:        "UNAUTHORIZED",
	Forbidden:           "FORBIDDEN",
	TimeoutRequest:      "TIMEOUT_REQUEST",
}

type customValidationErrCode map[string]string

var CustomValidationErrCode = customValidationErrCode{}

type errorMessage struct {
	SqlxNoRow       string
	InvalidDataType string
	InvalidRequest  string
	BadCredential   string
	TokenExpired    string
}

var ErrorMessage = errorMessage{
	SqlxNoRow:       "sql: no rows in result set",
	InvalidDataType: "invalid data type",
	InvalidRequest:  "invalid request",
	BadCredential:   "bad credential",
	TokenExpired:    "token has invalid claims: token is expired",
}
