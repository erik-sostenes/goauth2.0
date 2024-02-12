package business

import (
	"context"

	"github.com/erik-sostenes/auth-api/internal/repository"
	"github.com/erik-sostenes/auth-api/internal/repository/persistence"
)

type Exchanger interface {
	Exchange(context.Context, string) (string, error)
}

type exchanger struct {
	exchanger repository.CodeExchanger
	saver     persistence.UserSaver
	asker     repository.UserInfoAsker
	generator TokenGenerator
}

func NewExchanger(
	codeExchanger repository.CodeExchanger,
	saver persistence.UserSaver,
	asker repository.UserInfoAsker,
	generator TokenGenerator,
) Exchanger {
	return &exchanger{
		exchanger: codeExchanger,
		saver:     saver,
		asker:     asker,
		generator: generator,
	}
}

func (e *exchanger) Exchange(ctx context.Context, code string) (_ string, err error) {
	token, err := e.exchanger.ExchangeCode(ctx, code)
	if err != nil {
		return
	}

	user, err := e.asker.AskUserInfo(ctx, token)
	if err != nil {
		return
	}

	// create user
	if err = e.saver.Save(ctx, &user); err != nil {
		return
	}

	// generate token
	// payload JWT
	// {
	//	"id": 1,
	// 	"email": "contact@example.com",
	// 	"iss": "goauth",
	//  "exp": 24hrs
	// }
	return e.generator.GenerateToken(user.ID(), user.Email())
}
