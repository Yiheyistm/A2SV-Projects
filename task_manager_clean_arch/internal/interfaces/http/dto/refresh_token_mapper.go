package dto

import "github.com/yiheyistm/task_manager/internal/domain"

// change to refresh token to domain model
func (r *RefreshTokenRequest) FromRequestToDomainRefreshToken() *domain.RefreshToken {
	return &domain.RefreshToken{
		RefreshToken: r.RefreshToken,
	}
}

func FromDomainRefreshTokenToResponse(refreshToken *domain.RefreshToken) *domain.RefreshToken {
	return &domain.RefreshToken{
		AccessToken:  refreshToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
	}
}
