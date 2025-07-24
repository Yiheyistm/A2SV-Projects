// dto/user_mapper.go
package dto

import "github.com/yiheyistm/task_manager/internal/domain"

func (r *UserRequest) FromRequestToDomainUser() *domain.User {
	return &domain.User{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
		Role:     r.Role,
	}
}

func FromDomainUserToResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
}
func FromDomainUserToResponseList(users []domain.User) []UserResponse {
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *FromDomainUserToResponse(&user))
	}
	return userResponses
}
