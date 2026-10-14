package service

import (
	"banking-auth/domain"
	"banking-auth/dto"
	"banking-auth/logger"
	"fmt"

	"github.com/LzKelvinSoares/lzk-go-lib/errs"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(dto.LoginRequest) (*string, *errs.AppError)
	Verify(map[string]string) *errs.AppError
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

func (s DefaultAuthService) Login(req dto.LoginRequest) (*string, *errs.AppError) {
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

func (s DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	jwtToken, err := jwtTokenFromString(urlParams["token"])
	if err != nil {
		return err
	}

	if jwtToken.Valid {
		claims, ok := jwtToken.Claims.(*domain.Claims)
		if !ok {
			return errs.NewAuthorizationError("invalid token claims")
		}

		if claims.IsUserRole() {
			if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
				return errs.NewAuthorizationError("Request not verified with the token claims")
			}
		}

		isAuthorized := s.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
		if !isAuthorized {
			return errs.NewAuthorizationError(fmt.Sprintf("%s role is not authorized", claims.Role))
		}
		return nil
	} else {
		return errs.NewAuthorizationError("Invalid Token")
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, *errs.AppError) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})

	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return nil, errs.NewAuthorizationError(err.Error())
	}
	return token, nil
}

func NewDefaultAuthService(repo domain.AuthRepository, rolePermissions domain.RolePermissions) DefaultAuthService {
	return DefaultAuthService{
		repo:            repo,
		rolePermissions: rolePermissions,
	}
}
