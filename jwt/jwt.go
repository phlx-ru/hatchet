package jwt

import (
	"time"

	jwtv4 "github.com/golang-jwt/jwt/v4"
)

func Check(secret string) func(token *jwtv4.Token) (interface{}, error) {
	return func(token *jwtv4.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}

func Make(issuer string, secret string) string {
	claims := &jwtv4.RegisteredClaims{
		ExpiresAt: jwtv4.NewNumericDate(time.Now().Add(365 * 24 * time.Hour)),
		Issuer:    issuer,
	}

	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return signedString
}
