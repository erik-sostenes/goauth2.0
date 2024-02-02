package models

import (
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
	var userErr UserError

	if strings.TrimSpace(id) == "" {
		userErr |= MissingUserID
	}

	if strings.TrimSpace(name) == "" {
		userErr |= MissingUserName
	}

	// TODO: validate email
	if strings.TrimSpace(email) == "" {
		userErr |= MissingUserEmail
	}

	// TODO: validate profile url
	if strings.TrimSpace(picture) == "" {
		userErr |= MissingUserPicture
	}

	if userErr > 0 {
		return User{}, userErr
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
