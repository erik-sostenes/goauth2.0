package repository

import (
	"context"
	"fmt"

	"github.com/erik-sostenes/auth-api/internal/models"

	"golang.org/x/oauth2"
)

type (
	CodeExchanger interface {
		ExchangeCode(context.Context, string) (models.Token, error)
	}
)

type codeExchanger struct {
	config *oauth2.Config
}

func NewCodeExchanger(config *oauth2.Config) CodeExchanger {
	return &codeExchanger{
		config: config,
	}
}

func (c codeExchanger) ExchangeCode(ctx context.Context, code string) (models.Token, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return token, fmt.Errorf("%w: %s", models.ErrDuringCodeExchanger, err.Error())
	}

	return token, err
}

type codeExchangerMock struct {
	config *oauth2.Config
}

func NewCodeExchangerMock(config *oauth2.Config) CodeExchanger {
	return &codeExchangerMock{
		config: config,
	}
}

func (c *codeExchangerMock) ExchangeCode(ctx context.Context, code string) (models.Token, error) {
	return &oauth2.Token{}, nil
}
