package memory

import (
	"context"
	"fmt"

	"github.com/erik-sostenes/auth-api/internal/models"
	"github.com/erik-sostenes/auth-api/pkg"
)

type userMemory struct {
	// TODO: change value of primitive string by VO(models.UserId)
	set *pkg.Set[string, *models.User]
}

func NewUserMemory(set *pkg.Set[string, *models.User]) *userMemory {
	if set == nil {
		panic("missing set dependency")
	}

	return &userMemory{
		set: set,
	}
}

func (u *userMemory) Save(_ context.Context, user *models.User) error {
	if exists := u.set.Exist(user.ID()); !exists {
		u.set.Add(user.ID(), user)
	}

	return nil
}

func (u *userMemory) Get(_ context.Context, userID string) (*models.User, error) {
	user := u.set.GetByItem(userID)
	if user == nil {
		return nil, fmt.Errorf("%w: user with id '%s' not fount", models.UserNotFound, userID)
	}

	return user, nil
}
