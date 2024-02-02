package business

import (
	"context"
	"log"

	"github.com/erik-sostenes/auth-api/internal/repository"
)

type Exchanger interface {
	Exchange(context.Context, string) error
}

type exchanger struct {
	exchanger repository.CodeExchanger
	asker     repository.UserInfoAsker
}

func NewExchanger(codeExchanger repository.CodeExchanger, asker repository.UserInfoAsker) Exchanger {
	return &exchanger{
		exchanger: codeExchanger,
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

	// TODO: create user if it does not exist and generate JWT
	//
	// payload JWT
	// {
	//	"id": 1,
	// 	"email": "contact@example.com",
	// 	"iss": "crumbs",
	//  "exp": 24hrs
	// }

	log.Println(user)

	return
}
