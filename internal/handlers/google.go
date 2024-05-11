package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/erik-sostenes/auth-api/internal/business"
	"github.com/erik-sostenes/auth-api/internal/handlers/api"

	"github.com/labstack/echo/v4"
)

const cookieName = "state"

type GoogleLoginOAuthHandler interface {
	Login(echo.Context) error
}

type googleLoginOAuthHandler struct {
	provider business.PageProvider
}

func NewGoogleLoginOAuthHandler(provider business.PageProvider) GoogleLoginOAuthHandler {
	return &googleLoginOAuthHandler{
		provider: provider,
	}
}

func (g *googleLoginOAuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	page, state, err := g.provider.ProvidePage(ctx)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    state,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}

	apierr := api.SetCookie(c, cookie)
	if apierr != nil {
		return apierr
	}

	return c.HTML(http.StatusOK, page.String())
}

type GoogleCallbackOAuthHandler interface {
	Callback(echo.Context) error
}

type googleCallbackOAuthHandler struct {
	exchanger business.Exchanger
}

func NewGoogleCallbackOAuthHandler(exchanger business.Exchanger) GoogleCallbackOAuthHandler {
	return &googleCallbackOAuthHandler{
		exchanger: exchanger,
	}
}

func (g *googleCallbackOAuthHandler) Callback(c echo.Context) error {
	ctx := c.Request().Context()
	state := c.QueryParam(cookieName)

	cookie, apierr := api.GetCookie(c, cookieName)
	if apierr != nil {
		return apierr
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return err
	}

	if state != string(value) {
		return fmt.Errorf("%w", api.InvalidState)
	}

	code := c.QueryParam("code")
	token, err := g.exchanger.Exchange(ctx, code)
	if err != nil {
		return err
	}

	const redirectUrl = "https://www.google.com/"
	path, err := url.Parse(redirectUrl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	params := url.Values{}
	params.Set("params:Authorization:Bearer", token)
	path.RawQuery = params.Encode()

	return c.Redirect(http.StatusSeeOther, path.String())
}
