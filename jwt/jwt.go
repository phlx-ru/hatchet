package jwt

import (
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func Check(secret string) func(token *jwtv5.Token) (interface{}, error) {
	return func(_ *jwtv5.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}

func Make(issuer string, secret string) string {
	claims := &jwtv5.RegisteredClaims{
		ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(365 * 24 * time.Hour)),
		Issuer:    issuer,
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return signedString
}
