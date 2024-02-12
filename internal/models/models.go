package models

import (
	"fmt"
	"strings"

	"golang.org/x/oauth2"
)

type (
	Token *oauth2.Token

	User struct {
		id            string
		name          string
		email         string
		picture       string
		verifiedEmail bool
	}
)

func NewUser(id, name, email, picture string, verifiedEmail bool) (User, error) {
	if strings.TrimSpace(id) == "" {
		return User{}, fmt.Errorf("%w: user id '%s' is invalid or empty", MissingUserID, id)
	}

	if strings.TrimSpace(name) == "" {
		return User{}, fmt.Errorf("%w: user name '%s' is invalid or empty", MissingUserName, name)
	}

	// TODO: validate email
	if strings.TrimSpace(email) == "" {
		return User{}, fmt.Errorf("%w: user email '%s' is invalid or empty", MissingUserEmail, email)
	}

	// TODO: validate profile url
	if strings.TrimSpace(picture) == "" {
		return User{}, fmt.Errorf("%w: user picture '%s' is invalid or empty", MissingUserEmail, picture)
	}

	return User{
		id:            id,
		name:          name,
		email:         email,
		picture:       picture,
		verifiedEmail: verifiedEmail,
	}, nil
}

func (u User) ID() string {
	return u.id
}

func (u User) Name() string {
	return u.name
}

func (u User) Email() string {
	return u.email
}

func (u User) Picture() string {
	return u.picture
}

func (u User) VerifiedEmail() bool {
	return u.verifiedEmail
}
