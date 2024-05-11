package business

import (
	"context"
	"errors"

	"github.com/erik-sostenes/auth-api/internal/models"
	"github.com/erik-sostenes/auth-api/internal/repository"
	"github.com/erik-sostenes/auth-api/internal/repository/persistence"
)

type Exchanger interface {
	// Exchange method that exchanges the code for the token and obtains the user's data information to generate the token using standard jwt
	Exchange(context.Context, string) (string, error)
}

type exchanger struct {
	exchanger repository.CodeExchanger
	getter    persistence.UserGetter
	saver     persistence.UserSaver
	asker     repository.UserInfoAsker
	generator TokenGenerator
}

func NewExchanger(
	codeExchanger repository.CodeExchanger,
	getter persistence.UserGetter,
	saver persistence.UserSaver,
	asker repository.UserInfoAsker,
	generator TokenGenerator,
) Exchanger {
	if codeExchanger == nil {
		panic("missing codeExchanger dependency")
	}

	if getter == nil {
		panic("missing getter dependency")
	}

	if saver == nil {
		panic("missing saver dependency")
	}

	if asker == nil {
		panic("missing asker dependency")
	}

	if generator == nil {
		panic("missing generator dependency")
	}

	return &exchanger{
		exchanger: codeExchanger,
		getter:    getter,
		saver:     saver,
		asker:     asker,
		generator: generator,
	}
}

func (e *exchanger) Exchange(ctx context.Context, code string) (string, error) {
	token, err := e.exchanger.ExchangeCode(ctx, code)
	if err != nil {
		return "", err
	}

	user, err := e.asker.AskUserInfo(ctx, token)
	if err != nil {
		return "", err
	}

	// generate token
	// payload JWT
	// {
	//	"id": 1,
	// 	"email": "contact@example.com",
	// 	"iss": "goauth",
	//  "exp": 24hrs
	// }
	JWToken, err := e.generator.GenerateToken(user.ID(), user.Email())
	if err != nil {
		return "", err
	}

	// validate that the user exits in the DB
	_, err = e.getter.Get(ctx, user.ID())
	if err != nil {
		switch {
		case errors.Is(err, models.UserNotFound):
			// create user
			if err = e.saver.Save(ctx, user); err != nil {
				return "", err
			}
			return JWToken, nil
		default:
			return "", err
		}
	}

	if err = e.saver.Save(ctx, user); err != nil {
		return "", err
	}
	return JWToken, nil
}
