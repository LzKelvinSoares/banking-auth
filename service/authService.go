package service

import (
	"banking-auth/domain"
	"banking-auth/dto"
	"fmt"

	"github.com/LzKelvinSoares/lzk-go-lib/errs"
	"github.com/LzKelvinSoares/lzk-go-lib/logger"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(map[string]string) *errs.AppError
	Refresh(refreshTokenRequest dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError)
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

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	var appErr *errs.AppError
	var login *domain.Login
	if login, appErr = s.repo.FindBy(req.Username, req.Password); appErr != nil {
		return nil, appErr
	}

	claims := login.ClaimsForAccessToken()

	authToken := domain.NewAuthToken(claims)
	var accessToken, refreshToken string
	if accessToken, appErr = authToken.NewAccessToken(); appErr != nil {
		return nil, appErr
	}

	if refreshToken, appErr = s.repo.GenerateAndSaveRefreshTokenToStore(authToken); appErr != nil {
		return nil, appErr
	}
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	jwtToken, err := jwtTokenFromString(urlParams["token"])
	if err != nil {
		return err
	}

	if jwtToken.Valid {
		claims, ok := jwtToken.Claims.(*domain.AccessTokenClaims)
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

func (a DefaultAuthService) Refresh(request dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError) {
	if err := request.AccessTokenValid(); err != nil {
		return nil, errs.NewAuthenticationError("Invalid Token")
	}

	var appError *errs.AppError
	if appError = a.repo.RefreshTokenExist(request.RefreshToken); appError != nil {
		return nil, appError
	}

	var access_token string
	if access_token, appError = domain.NewAccessTokenFromRefreshToken(request.RefreshToken); appError != nil {
		return nil, appError
	}
	return &dto.LoginResponse{
		AccessToken: access_token,
	}, nil
}

func jwtTokenFromString(tokenString string) (*jwt.Token, *errs.AppError) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
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
