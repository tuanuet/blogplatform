package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import "context"

type SocialUserInfo struct {
	ProviderID string
	Name       string
	Email      string
}

type SocialAuthService interface {
	GetUserInfo(ctx context.Context, provider string, code string) (*SocialUserInfo, error)
	GetAuthURL(ctx context.Context, provider string) (string, error)
}

type socialAuthService struct{}

func NewSocialAuthService() SocialAuthService {
	return &socialAuthService{}
}

func (s *socialAuthService) GetUserInfo(ctx context.Context, provider string, code string) (*SocialUserInfo, error) {
	return nil, nil
}

func (s *socialAuthService) GetAuthURL(ctx context.Context, provider string) (string, error) {
	return "", nil
}
