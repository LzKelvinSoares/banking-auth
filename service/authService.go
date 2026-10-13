package service

import (
	"banking-auth/domain"
	"banking-auth/dto"
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(dto.LoginRequest) (*string, error)
	Verify(map[string]string) (bool, error)
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

type authTokenClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func (c authTokenClaims) IsUserRole() bool {
	return c.Role == "USER" || c.Role == "user"
}

func (c authTokenClaims) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) bool {
	if userID, ok := urlParams["userId"]; ok && userID != "" {
		return c.Subject == userID
	}
	return true
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*string, error) {
	login, err := s.repo.FindBy(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	token, tokErr := login.GenerateToken()
	if tokErr != nil {
		return nil, tokErr
	}

	return &token, nil
}

func (s DefaultAuthService) Verify(urlParams map[string]string) (bool, error) {
	jwtToken, err := jwtTokenFromString(urlParams["token"])
	if err != nil {
		return false, err
	}

	if !jwtToken.Valid {
		return false, errors.New("Invalid Token")
	}

	claims, ok := jwtToken.Claims.(*authTokenClaims)
	if !ok {
		return false, errors.New("invalid token claims")
	}

	if claims.IsUserRole() {
		if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
			return false, nil
		}
	}

	isAuthorized := s.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
	return isAuthorized, nil
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})

	if err != nil {
		log.Println("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

func NewDefaultAuthService(repo domain.AuthRepository, rolePermissions domain.RolePermissions) DefaultAuthService {
	return DefaultAuthService{
		repo:            repo,
		rolePermissions: rolePermissions,
	}
}
