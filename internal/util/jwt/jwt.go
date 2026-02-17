package jwt

import (
	"errors"
	"strings"
	"time"

	core "github.com/dgrijalva/jwt-go/v4"
)

func NewJwtClaims(subject, domain, role string, expired time.Duration) *JwtClaims {
	now := time.Now()
	expiresAt := now.Add(expired * time.Second)
	return &JwtClaims{
		StandardClaims: core.StandardClaims{
			IssuedAt:  core.At(now),
			ExpiresAt: core.At(expiresAt),
			NotBefore: core.At(now),
			Subject:   subject,
		},
		Domain: domain,
		Role:   role,
	}
}

var (
	DefaultSigningMethod = core.SigningMethodHS512
	BearerPrefix         = "Bearer "
)

type (
	TokenExpiredError = core.TokenExpiredError
)

func IsTokenExpiredError(err error) bool {
	var expErr *TokenExpiredError
	return errors.As(err, &expErr)
}

type JwtClaims struct {
	core.StandardClaims
	Domain, Role string
}

func GenerateToken(signingMethod core.SigningMethod, signingKey string, claims *JwtClaims) (string, error) {
	return core.NewWithClaims(signingMethod, claims).SignedString([]byte(signingKey))
}

func ParseToken(signingMethod core.SigningMethod, signingKey string, tokenString string) (*JwtClaims, error) {
	token, err := core.ParseWithClaims(strings.TrimPrefix(tokenString, BearerPrefix), &JwtClaims{}, func(t *core.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, &core.TokenNotValidYetError{}
	}
	return token.Claims.(*JwtClaims), nil
}
