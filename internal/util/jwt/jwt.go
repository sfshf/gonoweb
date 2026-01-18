package util

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

func NewJwtClaims(subject, domain, role string, expired time.Duration) *JwtClaims {
	now := time.Now()
	expiresAt := now.Add(expired * time.Second)
	return &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  jwt.At(now),
			ExpiresAt: jwt.At(expiresAt),
			NotBefore: jwt.At(now),
			Subject:   subject,
		},
		Domain: domain,
		Role:   role,
	}
}

var (
	DefaultSigningMethod = jwt.SigningMethodHS512
	BearerPrefix         = "Bearer "
)

type (
	TokenExpiredError = jwt.TokenExpiredError
)

func IsTokenExpiredError(err error) bool {
	var expErr *TokenExpiredError
	return errors.As(err, &expErr)
}

type JwtClaims struct {
	jwt.StandardClaims
	Domain, Role string
}

func GenerateToken(signingMethod jwt.SigningMethod, signingKey string, claims *JwtClaims) (string, error) {
	return jwt.NewWithClaims(signingMethod, claims).SignedString([]byte(signingKey))
}

func ParseToken(signingMethod jwt.SigningMethod, signingKey string, tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(strings.TrimPrefix(tokenString, BearerPrefix), &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, &jwt.TokenNotValidYetError{}
	}
	return token.Claims.(*JwtClaims), nil
}
