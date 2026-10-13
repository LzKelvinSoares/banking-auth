package service

import (
	"banking-auth/domain"
	"banking-auth/dto"
)

type AuthService interface {
	Login(dto.LoginRequest) (*string, error)
	// Verify(map[string]string) (bool, error)
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
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

// func (s DefaultAuthService) Verify(urlParams map[string]string) (bool, error) {

// }

func NewDefaultAuthService(repo domain.AuthRepository) DefaultAuthService {
	return DefaultAuthService{
		repo: repo,
	}
}
