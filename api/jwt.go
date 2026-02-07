package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type Claims struct {
	Payload
	jwt.RegisteredClaims
}

func LoadPrivateKey() *rsa.PrivateKey {
	data, err := os.ReadFile("keys/private.pem")
	if err != nil {
		panic(err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		panic(err)
	}
	return key
}

func LoadPublicKey() *rsa.PublicKey {
	data, err := os.ReadFile("keys/public.pem")
	if err != nil {
		panic(err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		panic(err)
	}
	return key
}

func GenJWT(private *rsa.PrivateKey, payload Payload) (string, error) {
	ActiveId := "340a2a39-d309-4d3c-b839-83df8a663e4c"

	claims := Claims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "http://localhost:8080",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
			Subject:   payload.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = ActiveId

	return token.SignedString(private)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	publicKey := LoadPublicKey()

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return publicKey, nil
		},
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer("http://localhost:8080"),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
