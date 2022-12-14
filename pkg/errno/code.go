package errno

const (
	SuccessCode = 0
	// service error
	ServiceErrCode = 10001
	// General incoming parameter error
	ParamErrCode = 10101
	// User-related incoming parameter error
	LoginErrCode            = 10202
	UserNotExistErrCode     = 10203
	UserAlreadyExistErrCode = 10204
	TokenExpiredErrCode     = 10205
	TokenValidationErrCode  = 10206
	TokenInvalidErrCode     = 10207
)

var (
	Success             = NewErr(SuccessCode, "Success")
	ServiceErr          = NewErr(ServiceErrCode, "Service is unable to start successfully")
	ParamErr            = NewErr(ParamErrCode, "Wrong Parameter has been given")
	LoginErr            = NewErr(LoginErrCode, "Wrong username or password")
	UserNotExistErr     = NewErr(UserNotExistErrCode, "User does not exists")
	UserAlreadyExistErr = NewErr(UserAlreadyExistErrCode, "User already exists")
)
