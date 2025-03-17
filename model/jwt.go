package model

import "github.com/golang-jwt/jwt/v5"

type JwtPayloadClaim struct {
	jwt.RegisteredClaims
	UserId string `json:"userId"`
	Role   string `json:"role"`
}
