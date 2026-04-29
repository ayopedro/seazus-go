package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenValidator interface {
	Validate(token string) (*JWTClaims, error)
}

type jwtValidator struct {
	secret []byte
}

type JWTClaims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

func NewJWTValidator(secret string) TokenValidator {
	return &jwtValidator{
		secret: []byte(secret),
	}
}

func generateToken(secret string, userID string, ttl time.Duration) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (v *jwtValidator) Validate(tokenString string) (*JWTClaims, error) {
	var claims JWTClaims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return v.secret, nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &claims, nil
}
