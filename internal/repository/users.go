// Package repository package is the persistence layer that stores the data
package repository

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/erik-sostenes/auth-api/internal/models"
)

type UserInfoAsker interface {
	AskUserInfo(context.Context, models.Token) (*models.User, error)
}

type userInfoAsker struct{}

func NewUserInfoAsker() UserInfoAsker {
	return &userInfoAsker{}
}

func (userInfoAsker) AskUserInfo(ctx context.Context, token models.Token) (user *models.User, err error) {
	const endpoint = "https://www.googleapis.com/oauth2/v2/userinfo"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return
	}

	// Setting query params
	params := url.Values{}
	params.Set("access_token", token.AccessToken)
	request.URL.RawQuery = params.Encode()

	// Doing request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	// Getting user data
	userResponse := User{}
	err = json.NewDecoder(response.Body).Decode(&userResponse)
	if err != nil {
		return
	}

	return models.NewUser(
		userResponse.ID,
		userResponse.Name,
		userResponse.Email,
		userResponse.Picture,
		userResponse.VerifiedEmail,
	)
}
