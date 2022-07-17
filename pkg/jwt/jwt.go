package jwt

import (
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/golang-jwt/jwt"
)

//JWT signing Key
type JWT struct {
	SigningKey []byte
}

var (
	ErrTokenExpired     = errno.WithCode(errno.TokenExpiredErrCode, "Token expired")
	ErrTokenNotValidYet = errno.WithCode(errno.TokenValidationErrCode, "Token is not active yet")
	ErrTokenMalformed   = errno.WithCode(errno.TokenInvalidErrCode, "That's not even a token")
	ErrTokenInvalid     = errno.WithCode(errno.TokenInvalidErrCode, "Couldn't handle this token")
)

//Structured version of Claims Section, as referenced at https://tools.ietf.org/html/rfc7519#section-4.1 See examples for how to use this with your own claim types
type CustomClaims struct {
	Id          int64
	AuthorityId int64
	jwt.StandardClaims
}

func NewJWT(SigningKey []byte) *JWT {
	return &JWT{
		SigningKey,
	}
}

//CreateToken create a new token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//ParseToken parses the token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

//checkToken get userId by token
func (j *JWT) CheckToken(token string) (int64, error) {
	if token == "" {
		return -1, nil
	}
	claim, err := j.ParseToken(token)
	if err != nil {
		return -1, ErrTokenInvalid
	}
	return claim.Id, nil
}
