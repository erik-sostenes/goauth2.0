package persistence

import (
	"context"
	"fmt"

	"github.com/erik-sostenes/auth-api/internal/models"
)

type UserSaver interface {
	Save(context.Context, *models.User) error
}

type userSaver struct {
	// TODO: change value of primitive string by VO(models.UserId)
	*Set[string, *models.User]
}

func NewUserSaver(set *Set[string, *models.User]) UserSaver {
	return &userSaver{
		Set: set,
	}
}

func (u *userSaver) Save(ctx context.Context, user *models.User) error {
	if exists := u.Set.Exist(user.ID()); exists {
		return fmt.Errorf("%w, %s", models.UserAlreadyExists, "the user with id"+user.ID()+"already exists")
	}

	u.Set.Add(user.ID(), user)

	return nil
}
