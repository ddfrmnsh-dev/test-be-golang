package usecase

import (
	"api-backend-go/model"
	"api-backend-go/service"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticationUseCase interface {
	LoginUser(email, password string) (string, model.User, error)
}
type authUseCaseImpl struct {
	userUseCase UserUseCase
	jwtService  service.JwtService
}

func NewAuthenticationImpl(uc UserUseCase, jwtService service.JwtService) AuthenticationUseCase {
	return &authUseCaseImpl{
		userUseCase: uc,
		jwtService:  jwtService,
	}
}

func (au *authUseCaseImpl) LoginUser(email, password string) (string, model.User, error) {
	var user model.User
	var err error

	user, err = au.userUseCase.FindUserByEmail(email)
	if err != nil {
		return "", user, err
	}

	hasValidRole := false
	if user.Role == "admin" || user.Role == "guest" {
		hasValidRole = true
	}

	if !hasValidRole {
		return "", user, errors.New("you are not authorized to login as admin")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", user, errors.New("invalid credentias")
	}

	token, err := au.jwtService.CreateToken(user)
	if err != nil {
		return "", user, err
	}
	return token, user, nil
}
