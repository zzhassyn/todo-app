package auth_transport_http

import "github.com/zzhassyn/todo-app/internal/core/domain"

// UserDTOResponse is intentionally duplicated from users_transport_http's
// DTO of the same shape rather than imported: the auth feature should not
// depend on the users feature's transport package (only on its service,
// via the AuthService -> auth_service.UsersRegistry chain). Keeping the
// DTO local preserves feature independence.
type UserDTOResponse struct {
	ID          int     `json:"id"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Email       string  `json:"email"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}
}
