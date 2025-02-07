package common

import (
	"fmt"
	"strconv"
	"template/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JwtGenerate(user_id int) (string, error) {
	now := time.Now()
	registeredClaims := jwt.RegisteredClaims{
		Issuer:    config.JWT_ISSUER,
		IssuedAt:  &jwt.NumericDate{Time: now},
		ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Hour)},
		Subject:   strconv.Itoa(user_id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)

	tokenStr, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func JwtGetID(tokenStr string) (int, error) {
	token, err := jwt.Parse(
		tokenStr,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return -1, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.JWT_SECRET), nil
		},
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(config.JWT_ISSUER),
	)
	if err != nil {
		return -1, AuthenticationError{Message: "JWT parsing failed"}
	}

	tokenMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return -1, AuthenticationError{Message: "JWT claims invalid"}
	}

	id_str, err := tokenMap.GetSubject()
	if err != nil {
		return -1, AuthenticationError{Message: "JWT ID missing"}
	}

	id, err := strconv.Atoi(id_str)
	if err != nil {
		return -1, AuthenticationError{Message: "JWT ID invalid"}
	}

	return id, nil
}
