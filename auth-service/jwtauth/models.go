package jwtauth

import "github.com/golang-jwt/jwt/v5"

type claims struct {
	jwt.RegisteredClaims
}
