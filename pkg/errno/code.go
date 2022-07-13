package errno

const (
	SuccessCode = 0
	//service error
	ServiceErrCode    = 10001
	ParamParseErrCode = 10002
	//General incoming parameter error
	ParamErrCode = 10101
	//User-related incoming parameter error
	LoginErrCode              = 10202
	UserNotExistErrCode       = 10203
	UserAlreadyExistErrCode   = 10204
	TokenExpiredErrCode       = 10205
	TokenValidationErrCode    = 10206
	TokenInvalidErrCode       = 10207
	UserNameValidationErrCode = 10208
	PasswordValidationErrCode = 10209
	//Video-related incoming parameter error
	VideoDataGetErrCode  = 10301
	VideoDataCopyErrCode = 10302
	//Comment-related incoming parameter error
	CommentTextErrCode = 10401
	//Relation-related incoming parameter error
	ActionTypeErrCode = 10501
)

var (
	Success = NewErrNo(SuccessCode, "Success")
	//service error
	ServiceErr    = NewErrNo(ServiceErrCode, "Service is unable to start successfully")
	ParamParseErr = NewErrNo(ParamParseErrCode, "Could not parse the param")
	//General incoming parameter error
	ParamErr = NewErrNo(ParamErrCode, "Wrong Parameter has been given")
	//User-related incoming parameter error
	LoginErr              = NewErrNo(LoginErrCode, "Wrong username or password")
	UserNotExistErr       = NewErrNo(UserNotExistErrCode, "User does not exists")
	UserAlreadyExistErr   = NewErrNo(UserAlreadyExistErrCode, "User already exists")
	TokenExpiredErr       = NewErrNo(TokenExpiredErrCode, "Token has been expired")
	TokenValidationErr    = NewErrNo(TokenInvalidErrCode, "Token is not active yet")
	TokenInvalidErr       = NewErrNo(TokenInvalidErrCode, "Token Invalid")
	UserNameValidationErr = NewErrNo(UserNameValidationErrCode, "Username is invalid")
	PasswordValidationErr = NewErrNo(PasswordValidationErrCode, "Password is invalid")
	//Video-related incoming parameter error
	VideoDataGetErr  = NewErrNo(VideoDataGetErrCode, "Could not get video data")
	VideoDataCopyErr = NewErrNo(VideoDataCopyErrCode, "Could not copy video data")
	//Comment-related incoming parameter error
	CommentTextErr = NewErrNo(CommentTextErrCode, "Comment text too long")
	//Relation-related incoming parameter error
	ActionTypeErr = NewErrNo(ActionTypeErrCode, "Action type is invalid")
)
