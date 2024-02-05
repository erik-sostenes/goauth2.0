package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/erik-sostenes/auth-api/internal/business"
	"github.com/erik-sostenes/auth-api/internal/handlers/api"

	"github.com/labstack/echo/v4"
)

const cookieName = "state"

type GoogleOAuthHandler interface {
	Login(echo.Context) error
	Callback(echo.Context) error
}

type googleOAuthHandler struct {
	provider  business.PageProvider
	exchanger business.Exchanger
}

func NewGoogleOAuthHandler(provider business.PageProvider, exchanger business.Exchanger) GoogleOAuthHandler {
	return &googleOAuthHandler{
		provider:  provider,
		exchanger: exchanger,
	}
}

func (g *googleOAuthHandler) Login(c echo.Context) error {
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

	err = api.SetCookie(c, cookie)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// TODO: display page or redirect if user is logged in
	return c.HTML(http.StatusOK, page.String())
}

func (g *googleOAuthHandler) Callback(c echo.Context) error {
	ctx := c.Request().Context()
	state := c.QueryParam(cookieName)

	cookie, apierr := api.ReadCookie(c, cookieName)
	if apierr != nil {
		return echo.NewHTTPError(apierr.Status(), apierr.Error())
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return err
	}

	if state != string(value) {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	code := c.QueryParam("code")
	err = g.exchanger.Exchange(ctx, code)
	if err != nil {
		return err
	}
	// TODO: redirect url
	url := "https://www.google.com/"
	return c.Redirect(http.StatusPermanentRedirect, url)
}
