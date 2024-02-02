package repository

import (
	"context"

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
	return c.config.Exchange(ctx, code)
}
