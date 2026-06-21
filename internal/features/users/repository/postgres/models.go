package users_postgres_repository

import "github.com/zzhassyn/todo-app/internal/core/domain"

type UserModel struct {
	ID           int
	Version      int
	FullName     string
	PhoneNumber  *string
	Email        string
	PasswordHash string
}

func userDomainFromModel(user UserModel) domain.User {
	return domain.NewUser(
		user.ID,
		user.Version,
		user.FullName,
		user.PhoneNumber,
		user.Email,
		user.PasswordHash,
	)
}

func userDomainsFromModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))
	for i, user := range users {
		userDomains[i] = userDomainFromModel(user)
	}
	return userDomains
}
