package persistence

import "github.com/erik-sostenes/auth-api/internal/models"

type UserDTO struct {
	Id, Name, Email, Picture string
	VerifiedEmail            bool
}

func (u *UserDTO) ToDomainUser() (*models.User, error) {
	return models.NewUser(
		u.Id, u.Name, u.Email, u.Picture, u.VerifiedEmail,
	)
}
