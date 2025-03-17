package service

import (
	"api-backend-go/config"
	"api-backend-go/model"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	CreateToken(user model.User) (string, error)
	VerifyToken(token string) (model.JwtPayloadClaim, error)
}

type jwtService struct {
	cfg config.TokenConfig
}

func NewJwtService(cfg config.TokenConfig) JwtService {
	return &jwtService{cfg: cfg}
}

func (j *jwtService) CreateToken(user model.User) (string, error) {
	tokenKey := []byte(j.cfg.JwtSignatureKey)

	newId := strconv.Itoa(user.Id)

	// var roleNames []

	claims := model.JwtPayloadClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.AccessTokenLifeTime)),
		},
		UserId: newId,
		Role:   user.Role,
	}

	jwtNewClaim := jwt.NewWithClaims(j.cfg.JwtSigninMethod, claims)
	token, err := jwtNewClaim.SignedString(tokenKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *jwtService) VerifyToken(token string) (model.JwtPayloadClaim, error) {
	tokenParse, err := jwt.ParseWithClaims(token, &model.JwtPayloadClaim{}, func(t *jwt.Token) (interface{}, error) {
		return j.cfg.JwtSignatureKey, nil
	})

	if err != nil {
		return model.JwtPayloadClaim{}, err
	}

	claim, ok := tokenParse.Claims.(*model.JwtPayloadClaim)

	if !ok {
		return model.JwtPayloadClaim{}, fmt.Errorf("error claim")
	}

	return *claim, nil
}
