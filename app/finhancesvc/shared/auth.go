package shared

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

const signingMethod = "HS256"

var ErrInvalidAuthHeader = errors.New("Invalid authentication header")
var ErrInvalidAuthTokenSign = errors.New("Invalid authentication signing method")
var ErrInvalidAuthToken = errors.New("Invalid authentication token")

type UserJWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func ExtractAuthorizationToken(c *gin.Context) (string, error) {
	bearerToken := c.Request.Header.Get("Authorization")
	authHeaderParts := strings.Split(bearerToken, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		return "", ErrInvalidAuthHeader
	}

	return authHeaderParts[1], nil
}

func ExtractJWTClaims(tokenString string, jwtSecret []byte) (jwtClaims *UserJWTClaims, err error) {
	keyFunc := generateKeyFunc(jwtSecret)

	token, err := jwt.ParseWithClaims(tokenString, &UserJWTClaims{}, keyFunc)
	if err != nil {
		fmt.Println(err)
		return
	}

	if jwtClaims, ok := token.Claims.(*UserJWTClaims); ok && token.Valid {
		return jwtClaims, nil
	} else {
		return jwtClaims, err
	}
}

func generateKeyFunc(secret []byte) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != signingMethod {
			return nil, ErrInvalidAuthTokenSign
		}
		return secret, nil
	}
}
