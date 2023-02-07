package errno

const (
	SuccessCode = 0
	// service error
	ServiceErrCode = 10001
	// General incoming parameter error
	ParamErrCode = 10101
	// User-related incoming parameter error
	LoginErrCode              = 10202
	UserNotExistErrCode       = 10203
	UserAlreadyExistErrCode   = 10204
	TokenExpiredErrCode       = 10205
	TokenValidationErrCode    = 10206
	TokenInvalidErrCode       = 10207
	UserNameValidationErrCode = 10208
	PasswordValidationErrCode = 10209
	VideoDataGetErrCode       = 10301
	VideoDataCopyErrCode      = 10302
	// Comment-related incoming parameter error
	CommentTextErrCode = 10401
	// Relation-related incoming parameter error
	ActionTypeErrCode = 10501
)

var (
	Success               = NewErr(SuccessCode, "Success")
	ServiceErr            = NewErr(ServiceErrCode, "Service is unable to start successfully")
	ParamErr              = NewErr(ParamErrCode, "Wrong Parameter has been given")
	LoginErr              = NewErr(LoginErrCode, "Wrong username or password")
	UserNotExistErr       = NewErr(UserNotExistErrCode, "User does not exists")
	UserAlreadyExistErr   = NewErr(UserAlreadyExistErrCode, "User already exists")
	TokenExpiredErr       = NewErr(TokenExpiredErrCode, "Token has been expired")
	TokenValidationErr    = NewErr(TokenInvalidErrCode, "Token is not active yet")
	TokenInvalidErr       = NewErr(TokenInvalidErrCode, "Token Invalid")
	UserNameValidationErr = NewErr(UserNameValidationErrCode, "Username is invalid")
	PasswordValidationErr = NewErr(PasswordValidationErrCode, "Password is invalid")
	VideoDataGetErr       = NewErr(VideoDataGetErrCode, "Could not get video data")
	VideoDataCopyErr      = NewErr(VideoDataCopyErrCode, "Could not copy video data")
	ActionTypeErr         = NewErr(ActionTypeErrCode, "action type is invalid")
)
