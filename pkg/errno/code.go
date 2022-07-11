package errno

const (
	SuccessCode = 0
	//service error
	ServiceErrCode = 10001
	//General incoming parameter error
	ParamErrCode = 10101
	//User-related incoming parameter error
	LoginErrCode            = 10202
	UserNotExistErrCode     = 10203
	UserAlreadyExistErrCode = 10204
	TokenExpiredErrCode     = 10205
	TokenValidationErrCode  = 10206
	TokenInvalidErrCode     = 10207
)

var (
	Success = NewErrNo(SuccessCode, "Success")

	ServiceErr = NewErrNo(ServiceErrCode, "Service is unable to start successfully")

	ParamErr = NewErrNo(ParamErrCode, "Wrong Parameter has been given")

	LoginErr            = NewErrNo(LoginErrCode, "Wrong username or password")
	UserNotExistErr     = NewErrNo(UserNotExistErrCode, "User does not exists")
	UserAlreadyExistErr = NewErrNo(UserAlreadyExistErrCode, "User already exists")
	TokenExpiredErr     = NewErrNo(TokenExpiredErrCode, "Token has been expired")
	TokenValidationErr  = NewErrNo(TokenInvalidErrCode, "Token is not active yet")
	TokenInvalidErr     = NewErrNo(TokenInvalidErrCode, "Token Invalid")
)
