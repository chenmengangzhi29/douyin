package jwt

type JWT struct {
	SigningKey []byte
}

var (
	ErrTokenExpired = WithCode(TokenExpiredErrCode, "Token expired")
)
