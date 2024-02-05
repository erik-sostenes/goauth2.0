package business

import (
	"context"

	"github.com/erik-sostenes/auth-api/internal/repository"
	"github.com/erik-sostenes/auth-api/internal/repository/persistence"
)

type Exchanger interface {
	Exchange(context.Context, string) error
}

type exchanger struct {
	exchanger repository.CodeExchanger
	saver     persistence.UserSaver
	asker     repository.UserInfoAsker
}

func NewExchanger(
	codeExchanger repository.CodeExchanger,
	saver persistence.UserSaver,
	asker repository.UserInfoAsker,
) Exchanger {
	return &exchanger{
		exchanger: codeExchanger,
		saver:     saver,
		asker:     asker,
	}
}

func (e *exchanger) Exchange(ctx context.Context, code string) (err error) {
	token, err := e.exchanger.ExchangeCode(ctx, code)
	if err != nil {
		return
	}

	user, err := e.asker.AskUserInfo(ctx, token)
	if err != nil {
		return
	}

	if err = e.saver.Save(ctx, &user); err != nil {
		return
	}

	// TODO: create user if it does not exist and generate JWT
	//
	// payload JWT
	// {
	//	"id": 1,
	// 	"email": "contact@example.com",
	// 	"iss": "crumbs",
	//  "exp": 24hrs
	// }

	return
}
